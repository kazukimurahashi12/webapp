import React, { useState, ChangeEvent, FormEvent, useEffect, useLayoutEffect } from "react";
import { useRouter } from "next/router";
import Head from "next/head";
import * as Yup from "yup";
import { Box, Button, Grid, Stack, TextField, Typography } from "@mui/material";
import { useFormik } from "formik";
import axios from "axios";
import { GetServerSideProps } from "next";
import NextLink from "next/link";
import { Logo } from "../../../components/logo";

type Blog = {
  id: string;
  loginID: string;
  title: string;
  content: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date;
};

type BlogProps = {
  blog: Blog;
};

  export default function Edit({ edit }) {
  const router = useRouter();
  const [isEditing, setEditing] = useState(false);
  const [propsBlog, setBlogProps] = useState<Blog>({
    id: "",
    loginID: "",
    title: "",
    content: "",
    createdAt: new Date(),
    updatedAt: new Date(),
    deletedAt: new Date(),
  });

  const formik = useFormik({
    initialValues: {
      id: "id",
      title: "title",
      content: "content",
    },
    validationSchema: Yup.object({
      title: Yup.string().min(1).max(50).required("タイトルを入力してください"),
      content: Yup.string().min(1).max(8000).required("記事内容を入力してください"),
    }),
    onSubmit: async (values, helpers) => {},
  });

  useEffect(() => {
    const getBlogContent = async () => {
      const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";

      try {
        const response = await axios.get(`http://${hostname}:8080/blog/overview/post/${edit}`, {
          withCredentials: true,
        });
        
        const blog: Blog = response.data.blog;
        //レスポンス情報を設定
        setBlogProps(blog)
      } catch (error) {
        console.error("エラーが発生しました", error);
        if (error.response.status === 302) {
          console.error("ログインしていません。StatusCode:", error.response.status);
          router.push("/auth/login");
        }
      }
    };
    getBlogContent();
  }, []);

    useLayoutEffect(() => {
    if (propsBlog) {
      formik.setFieldValue("id", String(propsBlog.id)); //数値型からstringに変換
      formik.setFieldValue("loginID", propsBlog.loginID);
      formik.setFieldValue("title", propsBlog.title);
      formik.setFieldValue("content", propsBlog.content);
    }
  }, [propsBlog]);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    // ブログ投稿後にブログ一覧ページにリダイレクト
    router.push("/auth/overview");
  };

  const handlePost = async (event: React.MouseEvent<HTMLElement>) => {
    if (isEditing) {
      // 追加: 投稿処理が実行中の場合は何もしない
      return;
    }
    setEditing(true); // 追加: 投稿処理を実行中にセット
    console.log("通過B", JSON.stringify(formik.values));
    const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
    try {
      const response = await axios.post(`http://${hostname}:8080/blog/edit`, formik.values, {
        headers: {
          "Content-Type": "application/json", // JSON形式で送信するためのヘッダー設定
        },
        withCredentials: true,
      });
      console.log("ブログ記事を編集しました");
    } catch (error) {
      // エラー処理を行う場合はここに記述
      console.error("Error editing blog:", error);
    } finally {
      // 投稿処理が終了したので false にセット
      setEditing(false);
      router.push("/auth/overview");
    }
  };

  return (
    <>
      <Head>
        <title>Blog</title>
      </Head>
      <Box sx={{ p: 3 }}>
        <Box
          component={NextLink}
          href="/auth/overview"
          sx={{
            display: "inline-flex",
            height: 32,
            width: 32,
          }}
        >
          <Logo />
        </Box>
      </Box>
      <Box
        component="main"
        sx={{
          display: "flex",
          justifyContent: "center",
          py: 8,
        }}
      >
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Box sx={{ width: "100%" }}>
              <Typography variant="h1" component="h1" gutterBottom>
                Edit
              </Typography>
              <form onSubmit={handleSubmit}>
                <Stack spacing={2}>
                  <TextField
                    error={!!(formik.touched.title && formik.errors.title)}
                    helperText={formik.touched.title && formik.errors.title}
                    id="title"
                    label="Title"
                    onBlur={formik.handleBlur}
                    value={formik.values.title}
                    onChange={formik.handleChange}
                    fullWidth
                  />
                  <TextField
                    error={!!(formik.touched.content && formik.errors.content)}
                    helperText={formik.touched.content && formik.errors.content}
                    id="content"
                    label="Content"
                    onBlur={formik.handleBlur}
                    value={formik.values.content}
                    onChange={formik.handleChange}
                    multiline
                    rows={15}
                    fullWidth
                  />
                </Stack>
                <Button
                  size="medium"
                  sx={{ mt: 3 }}
                  type="submit"
                  variant="contained"
                  onClick={handlePost}
                  disabled={!formik.isValid}
                >
                  Submit
                </Button>
              </form>
            </Box>
          </Grid>
          <Grid item xs={12} md={6}>
            <Typography variant="h2" component="h2" gutterBottom>
              Preview
            </Typography>
            {formik.values.title && (
              <Box
                sx={{
                  width: "100%",
                  backgroundColor: "#d4d4d4",
                  padding: "16px",
                  borderRadius: "4px",
                }}
              >
                <Typography variant="body1" style={{ marginTop: "8px", whiteSpace: "pre-wrap" }}>
                  {formik.values.title}
                </Typography>
              </Box>
            )}
            {formik.values.content && (
              <Box
                sx={{
                  width: "100%",
                  backgroundColor: "#f5f5f1",
                  padding: "16px",
                  borderRadius: "4px",
                }}
              >
                <Typography variant="body1" style={{ whiteSpace: "pre-wrap" }}>
                  {formik.values.content}
                </Typography>
              </Box>
            )}
          </Grid>
        </Grid>
      </Box>
    </>
  );
};

export async function getServerSideProps(context) {
  // ここで context.params を使用して id を取得
  const { edit } = context.params;

  return {
    props: {
      edit,
    },
  };
};

// export const getServerSideProps: GetServerSideProps<BlogProps> = async (context) => {
//   const { edit } = context.params; // idを取得
//   const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
//   const response = await axios.get(`http://${hostname}:8080/blog/overview/post/${edit}`, {
//     headers: {
//       "Content-Type": "aapplication/json",
//     },
//     withCredentials: true,
//   });

//   const blog: Blog = response.data.blog;

//   return {
//     props: {
//       blog,
//     },
//   };
// };


