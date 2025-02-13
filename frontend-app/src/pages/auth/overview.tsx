import Head from "next/head";
import ArrowUpOnSquareIcon from "@heroicons/react/24/solid/ArrowUpOnSquareIcon";
import ArrowDownOnSquareIcon from "@heroicons/react/24/solid/ArrowDownOnSquareIcon";
import PlusIcon from "@heroicons/react/24/solid/PlusIcon";
import {
  Box,
  Button,
  Container,
  Pagination,
  Stack,
  SvgIcon,
  Typography,
  Unstable_Grid2 as Grid,
} from "@mui/material";
import Layout from "../../layouts/dashboard/layout";
import { BlogCard } from "../../sections/blog/blog-card";
import React, { useEffect, useState } from "react";
import axios from "axios";
import { useRouter } from "next/router";

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
  blogs: Blog[];
};

type ClickValue = {
  value: string;
};



const Overview = ({ blogs }: BlogProps, { value }: ClickValue) => {
  const router = useRouter();

  const handleAdd = () => {
    // ブログ作成画面へ遷移
    router.push("/blog/post");
   
  };

  const handleBlogClick = (id: string) => {
    // ブログ詳細画面へ遷移
    router.push(`/blog/${id}`);
  
  };

  // blogsが未定義またはnullの場合、空の配列を初期値として設定
  const [blogsList, setBlogsList] = useState<Blog[]>([]);

  useEffect(() => {
    // ブログ情報を取得する関数
    const fetchBlogs = async () => {
      const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
      try {
        const response = await axios.get(`http://${hostname}:8080/blog/overview`, {
          withCredentials: true,
        });

        const blogsInfo = response.data.blogs;
        if (blogsInfo == null) {
          setBlogsList([]); // レスポンスのデータが null の場合、空の配列を設定
        } else {
          const blogs = blogsInfo.map((item: any) => ({
            id: item.id,
            loginID: item.loginID,
            title: item.title,
            content: item.content,
            createdAt: item.createdAt,
            updatedAt: item.updatedAt,
          }));
          setBlogsList(blogs); // レスポンスのデータがある場合、データを設定
        }
      } catch (error) {
        console.error("エラーが発生しました", error);
        if (error.response.status === 302 || error.response.status === 400) {
          router.push("/auth/login");
        }
      }
    };

    // ブログ情報を取得する関数を呼び出す
    fetchBlogs();
  }, []); // 空の依存配列を渡すことで、初回のみ実行される

  return (
    <>
      <Head>
        <title>Blog</title>
      </Head>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          py: 8,
        }}
      >
        <Container maxWidth="xl">
          <Stack spacing={3}>
            <Stack direction="row" justifyContent="space-between" spacing={4}>
              <Stack spacing={1}>
                <Typography variant="h4">Blog</Typography>
                <Stack alignItems="center" direction="row" spacing={1}></Stack>
              </Stack>
              <div>
                <Button
                  startIcon={
                    <SvgIcon fontSize="small">
                      <PlusIcon />
                    </SvgIcon>
                  }
                  variant="contained"
                  onClick={handleAdd}
                >
                  Create New Blog
                </Button>
              </div>
            </Stack>
            <Grid container spacing={3}>
              {blogsList.map((blog) => (
                <Grid xs={12} md={6} lg={4} key={blog.id}>
                  <BlogCard blog={blog} clickValue={value} onClick={() => handleBlogClick(blog.id)} />
                </Grid>
              ))}
            </Grid>
          </Stack>
        </Container>
      </Box>
    </>
  );
};

Overview.getLayout = (page: any) => <Layout>{page}</Layout>;

export default Overview;
