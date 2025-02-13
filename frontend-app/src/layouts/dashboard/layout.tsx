import { useCallback, useEffect, useState } from "react";
import { usePathname } from "next/navigation";
import { styled } from "@mui/material/styles";
import { withAuthGuard } from "../../hocs/with-auth-guard";
import { SideNav } from "./side-nav";
import { TopNav } from "./top-nav";
import axios from "axios";
import { useRouter } from "next/router";

const SIDE_NAV_WIDTH = 280;

const LayoutRoot = styled("div")(({ theme }) => ({
  display: "flex",
  flex: "1 1 auto",
  maxWidth: "100%",
  [theme.breakpoints.up("lg")]: {
    paddingLeft: SIDE_NAV_WIDTH,
  },
}));

const LayoutContainer = styled("div")({
  display: "flex",
  flex: "1 1 auto",
  flexDirection: "column",
  width: "100%",
});

type LayoutProps = {
  children?: React.ReactNode;
};

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const pathname = usePathname();
  const [openNav, setOpenNav] = useState(false);
  const router = useRouter();

  const handlePathnameChange = useCallback(() => {
    if (openNav) {
      setOpenNav(false);
    }
  }, [openNav]);

  const [fetchedID, setFetchedLoginID] = useState<string>("");

  // useEffect(() => {
  //   handlePathnameChange();
  //   const fetchData = async () => {
  //     const hostname = process.env.NODE_ENV === 'production' ? 'server-app' : 'localhost';
  //     try {
  //       const response = await axios.get(`http://${hostname}:8080/api/login-id`, {
  //         headers: {
  //           "Content-Type": "application/x-www-form-urlencoded",
  //         },
  //         withCredentials: true,
  //       });
  //       const fetchedID = response.data.id || null;
  //       setFetchedLoginID(fetchedID); // 取得した loginID を状態に設定する
  //     } catch (error) {
  //       console.error(error);
  //     }
  //   };
  //   // コンポーネントのマウント時にリクエストを実行
  //   fetchData();
  // }, [pathname]);

  useEffect(() => {
    handlePathnameChange();
    const getTopIdInfo = async () => {
      const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";

      try {
        const response = await axios.get(`http://${hostname}:8080/api/login-id`, {
          withCredentials: true,
        });
        
        const fetchedID = response.data.id;
        setFetchedLoginID(fetchedID); // 取得した loginID を状態に設定する
        
      } catch (error) {
        console.error("エラーが発生しました", error);
        if (error.response.status === 302) {
          console.error("ログインしていません。StatusCode:", error.response.status);
          router.push("/auth/login");
        }
      }
    };
    getTopIdInfo();
  }, [pathname]);

  return (
    <>
      <TopNav onNavOpen={() => setOpenNav(true)} />
      <SideNav onClose={() => setOpenNav(false)} open={openNav} loginID={fetchedID} />
      <LayoutRoot>
        <LayoutContainer>{children}</LayoutContainer>
      </LayoutRoot>
    </>
  );
};

export default withAuthGuard(Layout);
