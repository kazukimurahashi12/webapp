import * as React from "react";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import Markdown from "./markdown";
import { useEffect, useState } from "react";
import closestIndexTo from "date-fns/esm/closestIndexTo";
import Box from "@mui/material/Box";

type Blog = {
  id: string;
  loginID: string;
  title: string;
  content: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date;
};

interface MainProps {
  posts: Blog[];
  id: any;
}




export default function Main(props: MainProps) {
  const { posts } = props;
  const { id } = props;
  const [selectedPost, setSelectedPost] = useState<Blog | null>(null);
  const [formattedDate, setformattedDate] = useState<String | null>(null);

  
  const formatDate = (date: any) => {
    const options: Intl.DateTimeFormatOptions = {
      year: 'numeric',
      month: 'long',
    };
    if (date != null){
    return  new Date(date).toLocaleDateString('en-US', options);
    }
  };

  

  // idが変更されたときに選択された投稿を更新
  useEffect(() => {    
    if (id && posts.length > 0) {
      // findIndexメソッドを使用してidを検索し、indexを取得
      const index = posts.findIndex((blog) => blog.id === id);
      setSelectedPost(posts[index]);
      setformattedDate(formatDate(posts[index].createdAt));
    } else if (posts.length > 0) {
      setSelectedPost(posts[0]);
      setformattedDate(formatDate(posts[0].createdAt));
    } else {
      // const newPost = Blog[0].title = "aaa";
      // setSelectedPost(newPost);
    }
  }, [id, posts]);

  if (!selectedPost) {
    // 選択された投稿が存在しない場合の処理
    return <div>No posts available</div>;
  }

  return (
    <Grid
      item
      xs={12}
      md={8}
      sx={{
        "& .markdown": {
          py: 3,
        },
      }}
    >
  <Box display="flex" justifyContent="space-between">
        <Typography variant="h5" display="inline">
          {selectedPost.title}
        </Typography>
        <Typography variant="h6" display="inline" color="textSecondary">
          {selectedPost.loginID} {formattedDate} {/* スペースを追加 */}
        </Typography>
      </Box>
      <Divider />
      <Markdown className="markdown" key={selectedPost.title}>
        {selectedPost.content}
      </Markdown>
    </Grid>
  );
}
