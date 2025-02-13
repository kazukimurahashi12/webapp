import React, { useEffect, useState } from "react";
import Container from "@mui/material/Container";
import { Box, Typography } from "@mui/material";
import Chip from "@mui/material/Chip";
import { useRouter } from "next/router";
import { GetServerSideProps } from "next";
import axios from "axios";
import format from "date-fns/format";
import NextLink from "next/link";
import { Logo } from "../../components/logo";

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

export default function Blog({ id }) {
  const router = useRouter();
  const [propsBlog, setBlogProps] = useState<Blog>({
    id: "",
    loginID: "",
    title: "",
    content: "",
    createdAt: new Date(),
    updatedAt: new Date(),
    deletedAt: new Date(),
  });

  const content = {
    date: `${format(new Date(propsBlog.createdAt), "yyyy/MM/dd")}`,
    "header-p1": `${propsBlog.title}`,
    name: `${propsBlog.loginID}`,
    paragraph1: `${propsBlog.content}`,
  };

  useEffect(() => {
    const getBlogContent = async () => {
      const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";

      try {
        const response = await axios.get(`http://${hostname}:8080/blog/overview/post/${id}`, {
          withCredentials: true,
        });
        
        const blog: Blog = response.data.blog;
        setBlogProps(blog);

      } catch (error) {
        console.error("エラーが発生しました", error);
        if (error.response.status === 302) {
          router.push("/auth/login");
        }
      }
    };
    getBlogContent();
  }, []);

  return (
    <section>
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
      <Container maxWidth="md">
        <Box py={10}>
          <Box textAlign="center" mb={5}>
            <Container maxWidth="sm">
              <Chip color="primary" label={content["date"]} />
              <Box my={4}>
                <Typography variant="h3" component="h2">
                  <Typography variant="h3" component="span" color="primary">
                    {content["header-p1"]}{" "}
                  </Typography>
                  <Typography variant="h3" component="span">
                    {content["header-p2"]}
                  </Typography>
                </Typography>
              </Box>
              <Box display="flex" justifyContent="center" alignItems="center">
                <Box ml={2} textAlign="left">
                  <Typography variant="subtitle1" component="h2" style={{ lineHeight: 1 }}>
                    {content["name"]}
                  </Typography>
                  <Typography variant="subtitle1" component="h3" color="textSecondary">
                    {content["job"]}
                  </Typography>
                </Box>
              </Box>
            </Container>
          </Box>
          <Box>
            <Typography
              variant="subtitle1"
              style={{ whiteSpace: "pre-wrap", marginBottom: "16px" }}
              color="textPrimary"
              paragraph={true}
            >
              {content["paragraph1"]}
            </Typography>
            <Box my={4}>
              <img src={content["image"]} alt="" style={{ maxWidth: "100%", borderRadius: "4px" }} />
            </Box>
          </Box>
        </Box>
      </Container>
    </section>
  );
}

export async function getServerSideProps(context) {
  const { id } = context.params;

  return {
    props: {
      id,
    },
  };
}
