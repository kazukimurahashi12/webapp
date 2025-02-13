import * as React from 'react';
import CssBaseline from '@mui/material/CssBaseline';
import Grid from '@mui/material/Grid';
import Container from '@mui/material/Container';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import Header from './view/header';
import MainFeaturedPost from './view/mainFeaturedPost';
import Main from './view/main';
import Sidebar from './view/sidebar';
import Footer from './view/footer';
import axios from 'axios';
import { useEffect, useState } from 'react';
import router from 'next/router';
import { Typography } from '@mui/material';
import { FeaturedPost } from './view/featuredPost';


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

const mainFeaturedPost = {
  title: 'Engineer technical blog',
  description:
    "Life is either a daring adventure or nothing at all",
  image: 'https://source.unsplash.com/random?wallpapers',
  imageText: 'main image description',
  linkText: '',
};


const sidebar = 
{
  title: 'About',
  description:
    '徒然なるままに　エンジニアの技術ブログ',
  archives: 
  [
    { title: 'June 2023', url: '/' },
    { title: 'February 2020', url: '/' },
    { title: 'January 2020', url: '/' },
    { title: 'November 1999', url: '/' },
    { title: 'October 1999', url: '/' },
    { title: 'September 1999', url: '/' },
    { title: 'August 1999', url: '' },
    { title: 'July 1999', url: '/' },
    { title: 'June 1999', url: '/' },
    { title: 'May 1999', url: '/' },
    { title: 'April 1999', url: '/' },
  ],
};

const linkStyle = {
  cursor: 'pointer', // マウスカーソルをポインターに変更
};

const defaultTheme = createTheme();

export default function Blog() {
  const [blogsList, setBlogsList] = useState<Blog[]>([]);
  const [id, setId] = useState("");
  // const { postid } = id;

  const handleValueChange = (postId) => {
    setId(postId);
  };

  const handleSeeAll = () => {
  //  router.push();
  };

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
  
        }
      }
    };
    // ブログ情報を取得する関数を呼び出す
    fetchBlogs();
  }, []);
  
  return (
    <ThemeProvider theme={defaultTheme}>
      <CssBaseline />
      <Container maxWidth="lg">
        <Header title="Blog" sections={[]}/>
        <main>
          <MainFeaturedPost post={mainFeaturedPost} />
          <Grid container spacing={4}>
          {blogsList.slice(0, 2).map((post) => (
            <FeaturedPost key={post.title} post={post} handleValueChange={handleValueChange}/>
          ))}
          </Grid>
          <Typography variant="h6" color="primary" gutterBottom onClick={handleSeeAll} style={linkStyle}>
        {"...see all"}
      </Typography>
          <Grid container spacing={5} sx={{ mt: 3 }}>
            <Main posts={blogsList} id={id}/>
            <Sidebar
              title={sidebar.title}
              description={sidebar.description}
              archives={sidebar.archives}
            />
          </Grid>
        </main>
        <Footer
        title="Footer"
        description=""
      />
      </Container>
    </ThemeProvider>
  );
}
