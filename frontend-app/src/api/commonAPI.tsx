import axios from "axios";
import { GetServerSideProps } from "next";

type Blog = {
  ID: string;
  LoginID: string;
  Title: string;
  Content: string;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: Date;
};
  type BlogProps = {
    blog: Blog;
  };
  

//TODO 必要であれば再度実装
// export const getServerSideProps: GetServerSideProps<BlogProps> = async (context) => {
  
//     const response = await axios.get("http://localhost:8080/", {
//       headers: {
//         "Content-Type": "application/x-www-form-urlencoded",
//       },
//       withCredentials: true,
//     });
  
//     const blogsInfo = response.data.blogs;
  
//     const blogs: Blog[] = [blogsInfo].map((item: any) => ({
//       ID: item.ID,
//       LoginID: item.LoginID,
//       title: item.Title,
//       content: item.Content,
//       createdAt: item.CreatedAt,
//       updatedAt: item.UpdatedAt
//       }));
  
//     return {
//       props: {
//         blogs,
//       },
//     };
//   };