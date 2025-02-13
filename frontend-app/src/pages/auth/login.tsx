import { useCallback, useEffect, useState } from "react";
import Head from "next/head";
import NextLink from "next/link";
import { useRouter } from "next/router";
import { useFormik } from "formik";
import * as Yup from "yup";
import {
  Alert,
  Box,
  Button,
  FormHelperText,
  Link,
  Stack,
  Tab,
  Tabs,
  TextField,
  Typography,
} from "@mui/material";
import { GetServerSideProps, NextPage } from "next";
import axios from "axios";

import { useAuth } from "../../hooks/use-auth";
import { Layout as AuthLayout } from "../../layouts/auth/layout";
import React from "react";

type User = {
  userId: string;
};

type HomeProps = {
  user: User;
};

const mode: string = "login"

const Page: NextPage<HomeProps> & { getLayout: (page: React.ReactNode) => React.ReactNode } = () => {
  const router = useRouter();
  const auth = useAuth();
  const [method, setMethod] = useState("userId");
  const [userId, setUserId] = useState("root");
  const [password, setPassword] = useState("root");

  useEffect(() => {
    const fetchData = async () => {
      const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";

      try {
        const response = await axios.get(`http://${hostname}:8080/login`, {
          withCredentials: true,
        });
        const data = response.data;
        // レスポンスデータの処理
      } catch (error) {
        console.error(error);
      }
    };
    // コンポーネントのマウント時にリクエストを実行
    fetchData();
  }, []);

  const handleLogin = async (event: React.MouseEvent<HTMLElement>) => {
    const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
    try {
      const response = await axios.post(
        `http://${hostname}:8080/login`,
        { userId, password },
        {
          headers: {
            "Content-Type": "application/json",
          },
          withCredentials: true,
        }
      );
      // ログイン成功時の処理
      // await auth.signIn(userId, password);
      router.push("/auth/overview");
      // ログイン成功時の追加の処理を追記する場合はここに記述する
    } catch (error) {
      // ログイン失敗時の処理
      console.error(error);
      // ログイン失敗時の追加の処理を追記する場合はここに記述する
    }
  };

  const formik = useFormik({
    initialValues: {
      userId: "root",
      password: "root",
      submit: null,
    },
    validationSchema: Yup.object({
      userId: Yup.string().max(20).required("IDを入力してください"),
      password: Yup.string().max(20).required("Passwordを入力してください"),
    }),
    onSubmit: async (values, helpers) => {},
  });

  const handleMethodChange = useCallback((event: React.SyntheticEvent, value: string) => {
    setMethod(value);
  }, []);

  return (
    <>
      <Head>
        <title>Login</title>
      </Head>
      <Box
        sx={{
          backgroundColor: "background.paper",
          flex: "1 1 auto",
          alignItems: "center",
          display: "flex",
          justifyContent: "center",
        }}
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
              <Typography variant="h4">Login</Typography>
              <Typography color="text.secondary" variant="body2">
                Don&apos;t have an account? &nbsp;
                <Link
                  component={NextLink}
                  href="/auth/register"
                  underline="hover"
                  variant="subtitle2"
                >
                  Register
                </Link>
              </Typography>
            </Stack>
            <Tabs onChange={handleMethodChange} sx={{ mb: 3 }} value={method}>
              <Tab label="Acoount" value="userId" />
            </Tabs>
            {method === "userId" && (
              <form noValidate onSubmit={formik.handleSubmit}>
                <Stack spacing={3}>
                  <TextField
                    error={!!(formik.touched.userId && formik.errors.userId)}
                    fullWidth
                    helperText={formik.touched.userId && formik.errors.userId}
                    label="ID"
                    name="userId"
                    onBlur={formik.handleBlur}
                    onChange={(event) => {
                      const input = event.target.value.trim();
                      formik.handleChange(event);
                      setUserId(input);
                    }}
                    type="text"
                    value={userId && formik.values.userId}
                  />
                  <TextField
                    error={!!(formik.touched.password && formik.errors.password)}
                    fullWidth
                    helperText={formik.touched.password && formik.errors.password}
                    label="Password"
                    name="password"
                    onBlur={formik.handleBlur}
                    onChange={(event) => {
                      const input = event.target.value.trim();
                      formik.handleChange(event);
                      setPassword(input);
                    }}
                    type="password"
                    value={password && formik.values.password}
                  />
                </Stack>
                <Button
                  fullWidth
                  size="large"
                  sx={{ mt: 3 }}
                  variant="contained"
                  onClick={handleLogin}
                  disabled={!formik.isValid}
                >
                  Continue
                </Button>
                <Alert color="info" severity="info" sx={{ mt: 3 }}>
                  <div>
                    You can use <b>root</b> and password <b>root</b> (rootでログイン出来ます。)
                  </div>
                </Alert>
              </form>
            )}
          </div>
        </Box>
      </Box>
    </>
  );
};

Page.getLayout = (page: React.ReactNode) => <AuthLayout>{page}</AuthLayout>;

export default Page;