import * as React from "react";
import Toolbar from "@mui/material/Toolbar";
import Button from "@mui/material/Button";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";
import router from "next/router";

interface HeaderProps {
  sections: ReadonlyArray<{
    title: string;
    url: string;
  }>;
  title: string;
}

export default function Header(props: HeaderProps) {
  const { title } = props;

  const handleSingUp = () => {
    router.push("/auth/register");
  };

  const handleLogin = () => {
    router.push("/auth/login");
  };

  return (
    <React.Fragment>
      <Toolbar sx={{ borderBottom: 1, borderColor: "divider" }}>
        <Typography
          component="h2"
          variant="h5"
          color="inherit"
          align="center"
          noWrap
          sx={{ flex: 1 }}
        >
          {title}
        </Typography>
        <IconButton></IconButton>
        <Button size="small" variant="contained" onClick={handleSingUp}>
          Sign up
        </Button>
        {"　"}
        <Button size="small" variant="contained" onClick={handleLogin}>
          Login
        </Button>
      </Toolbar>
      {/* <Toolbar
        component="nav"
        variant="dense"
        sx={{ justifyContent: 'space-between', overflowX: 'auto' }}
      >
      </Toolbar> */}
    </React.Fragment>
  );
}
