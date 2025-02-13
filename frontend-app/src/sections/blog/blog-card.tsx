import { useState } from 'react';
import {
  Avatar,
  Box,
  Card,
  CardContent,
  Divider,
  MenuItem,
  Select,
  Stack,
  SvgIcon,
  Typography,
} from '@mui/material';
import { format } from 'date-fns';
import { SelectChangeEvent } from '@mui/material/Select';  // 追加

// import ArrowDownOnSquareIcon from '@heroicons/react/24/solid/ArrowDownOnSquareIcon';
import ClockIcon from '@heroicons/react/24/solid/ClockIcon';
import router from 'next/router';
import DeleteDialog from './delete-dialog';

type Blog = {
  id: string;
  loginID: string;
  title: string;
  content: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date;
};

type ClickValue = {
  value: string;
};

type BlogCardProps = {
  blog: Blog;
  clickValue: string;
  onClick: (id: string, action: string) => void;
};

export const BlogCard: React.FC<BlogCardProps> = ({ blog }) => {
  const [selectedAction, setSelectedAction] = useState('');
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const handleOpen = () => {
    setShowDeleteDialog(true);
  };

  const handleClose = () => {
    setShowDeleteDialog(false);
  };

  const handleActionChange = (event: SelectChangeEvent<string>) => {
    const selectedValue = event.target.value as string;
    setSelectedAction(selectedValue);

    if (selectedValue === "edit") {
      router.push(`/blog/edit/${blog.id}`); 
    } else if (selectedValue === "delete") {
      setShowDeleteDialog(true);
    }
  };

  const handleContainerClick = (event: React.MouseEvent<HTMLDivElement>) => {
    // プルダウン以外の部分がクリックされた場合の処理
      router.push(`/blog/${blog.id}`); // ブログ記事の詳細ページに遷移
  };

  const handleMouseEnter = (event: React.MouseEvent<HTMLDivElement>) => {
    event.currentTarget.style.cursor = "pointer";
  };
  const handleMouseLeave = (event: React.MouseEvent<HTMLDivElement>) => {
    event.currentTarget.style.cursor = "default";
  };

  return (
    <div>
      <Card
        sx={{
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
        }}
      >
        <CardContent onClick={handleContainerClick} onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}>
          <Typography align="center" gutterBottom variant="h5">
            {blog.title}
          </Typography>
          <Typography align="center" variant="body1">
            {blog.content.split('\n')[0]}
          </Typography>
        </CardContent>
        <Box sx={{ flexGrow: 1 }} />
        <Divider />
        <Stack alignItems="center" direction="row" justifyContent="space-between" spacing={2} sx={{ p: 2 }}>
          <Stack alignItems="center" direction="row" spacing={1}>
            <SvgIcon color="action" fontSize="small">
              <ClockIcon />
            </SvgIcon>
            <Typography color="text.secondary" display="inline" variant="body2">
              {format(new Date(blog.createdAt), 'yyyy/MM/dd HH:mm')}
            </Typography>
          </Stack>
          <Stack alignItems="center" direction="row" spacing={1}>
            {/* <SvgIcon color="action" fontSize="small">
              <ArrowDownOnSquareIcon />
            </SvgIcon> */}
            {/* <Typography color="text.secondary" display="inline" variant="body2">
              {blog.downloads} Good
            </Typography> */}
          </Stack>
          <Select value={selectedAction} onChange={handleActionChange} variant="outlined" size="small">
            <MenuItem value="edit">Edit</MenuItem>
            <MenuItem value="delete">Delete</MenuItem>
          </Select>
        </Stack>
        {showDeleteDialog && (
      <DeleteDialog
        id={blog.id} // ブログ記事のIDを渡す
        // handleClose={() => setShowDeleteDialog(false)} // ダイアログを閉じるためのハンドラー
      />
    )}
      </Card>
    </div>
  );
};

export default BlogCard;
