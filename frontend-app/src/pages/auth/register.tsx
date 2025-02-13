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
import Layout from "../../layouts/dashboard/layout";

type User = {
  UserId: string;
};

type HomeProps = {
  user: User;
};

const Register: NextPage<HomeProps> & { getLayout: (page: React.ReactNode) => React.ReactNode } = () => {
  const router = useRouter();
  const auth = useAuth();
  const [method, setMethod] = useState("userId");
  const [userId, setUserId] = useState("root");
  const [password, setPassword] = useState("root");

  // useEffect(() => {
  //   const fetchData = async () => {
  //     const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
  //     try {
  //       const response = await axios.get(`http://${hostname}:8080/login`, {
  //         headers: {
  //           "Content-Type": "application/x-www-form-urlencoded",
  //         },
  //         withCredentials: true,
  //       });
  //       const data = response.data;
  //       // レスポンスデータの処理
  //     } catch (error) {
  //       console.error(error);
  //     }
  //   };
  //   // コンポーネントのマウント時にリクエストを実行
  //   fetchData();
  // }, []);

  const handleRegist = async (event: React.MouseEvent<HTMLElement>) => {
    const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
    try {
      const response = await axios.post(
        `http://${hostname}:8080/regist`,
        { userId, password },
        {
          headers: {
            "Content-Type": "application/json",
          },
          withCredentials: true,
        }
      );

      // ログイン成功時の処理
      await auth.signIn(userId, password);
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
      userId: "",
      password: "",
      submit: null,
    },
    validationSchema: Yup.object({
      userId: Yup.string().min(2).max(10).required("IDを入力してください"),
      password: Yup.string().min(2).max(10).required("Passwordを入力してください"),
    }),
    onSubmit: async (values, helpers) => {},
  });

  const handleMethodChange = useCallback((event: React.SyntheticEvent, value: string) => {
    setMethod(value);
  }, []);

  const handleSkip = useCallback(() => {
    auth.skip();
    router.push("/auth/overview");
  }, [auth, router]);

  return (
    <>
      <Head>
        <title>Register</title>
      </Head>
      <Box
      sx={{ m: -1.5 }}
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
            <Typography variant="h4">
                Register
              </Typography>
              <Typography
                color="text.secondary"
                variant="body2"
              >
                Already have an account?
                &nbsp;
                <Link
                  component={NextLink}
                  href="/auth/login"
                  underline="hover"
                  variant="subtitle2"
                >
                  Log in
                </Link>
              </Typography>
            </Stack>
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
                  size="medium"
                  sx={{ mt: 3 }}
                  variant="contained"
                  onClick={handleRegist}
                  disabled={!formik.isValid }
                >
                  Continue
                </Button>
              </form>
          </div>
        </Box>
      </Box>
    </>
  );
};

// export const getServerSideProps: GetServerSideProps<HomeProps> = async () => {
//   const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
//   const response = await axios.get(`http://${hostname}:8080/login`,
//   {
//     withCredentials: true,
//   });
//   const user = response.data;
//   return {
//     props: {
//       user,
//     },
//   };
// };


Register.getLayout = (page) => (
  <Layout>
    {page}
  </Layout>
);


export default Register;
