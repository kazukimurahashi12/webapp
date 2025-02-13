import { useCallback, useEffect, useState } from "react";
import Head from "next/head";
import NextLink from "next/link";
import { useRouter } from "next/router";
import * as Yup from "yup";
import {
  Box,
  Button,
  Stack,
  Tab,
  TextField,
  Typography,
} from "@mui/material";
import { GetServerSideProps, NextPage } from "next";
import axios from "axios";
import React from "react";
import Layout from "../../layouts/dashboard/layout";

type User = {
  userId: string;
  password: string;
};

type HomeProps = {
  user: User;
};

interface BlogForm {
  loginID: string;
  title: string;
  content: string;
}

const Logout: NextPage<HomeProps> & { getLayout: (page: React.ReactNode) => React.ReactNode } = ({
  user,
}) => {
  const router = useRouter();
  const [userForm, setUserForm] = useState<User>({
    userId: "",
    password: "",
  });
  const [errorMessage, setErrorMessage] = useState<string>("");

  useEffect(() => {
    const fetchData = async () => {
      const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
      try {
        const response = await axios.get(`http://${hostname}:8080/api/login-id`, {
          withCredentials: true,
        });
        setUserForm((prevBlogForm) => ({
          ...prevBlogForm,
          userId: response.data.id,
        }));
      } catch (error) {
        console.error(error);
      }
    };
    // コンポーネントのマウント時にリクエストを実行
    fetchData();
  }, []);

  const handleLogout = async (event: React.MouseEvent<HTMLElement>) => {
    const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
    try {
      const response = await axios.post(
        `http://${hostname}:8080/logout`,
        userForm,
        {
          headers: {
            "Content-Type": "application/json",
          },
          withCredentials: true,
        }
      );
      router.push("/auth/login");
    } catch (error) {
      // ログイン失敗時の処理
      console.error(error);
    }
  };

  return (
    <>
      <Head>
        <title>Logout</title>
      </Head>
      <Box
      >
        <Box
          sx={{
            maxWidth: 550,
            px: 3,
            py: "100px",
            width: "100%",
          }}
        >
          <div>
            <Stack spacing={1} sx={{ mb: 3 }}>
              <Typography variant="h4">Logout</Typography>
              <Typography color="text.secondary" variant="body2"></Typography>
            </Stack>
            <Tab label="Acoount" />
              <Stack spacing={3}>
                <Box sx={{ p: 2, backgroundColor: "#f5f5f5", borderRadius: 4 }}>
                  <Typography variant="caption">ID</Typography>
                  <Typography variant="h5">{userForm.userId}</Typography>
                </Box>
              </Stack>
              <Stack spacing={3} sx={{ mt: 4 }}>
            <TextField
              type="password"
              label="Password"
              value={userForm.password}
              onChange={(event) => {
                const inputId = event.target.value.trim();
                setUserForm({ ...userForm, password: inputId });
                if (inputId === "") {
                  setErrorMessage("Password cannot be empty");
                } else {
                  setErrorMessage("");
                }
              }}
              required
              fullWidth
              error={!!errorMessage}
              helperText={errorMessage}
            />
          </Stack>
              <Button
                size="medium"
                sx={{ mt: 3 }}
                variant="contained"
                value={userForm.userId}
                onClick={handleLogout}
                disabled={!!errorMessage}
              >
                Logout
              </Button>
          </div>
        </Box>
      </Box>
    </>
  );
};

Logout.getLayout = (page: React.ReactNode) => (
  <Layout>
    {page}
  </Layout>
);

export default Logout;
