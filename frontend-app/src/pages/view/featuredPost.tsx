import * as React from 'react';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';
import Card from '@mui/material/Card';
import CardActionArea from '@mui/material/CardActionArea';
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import router from 'next/router';
import Blog from '..';

type FeaturedPostProps = {
  handleValueChange(id: string): unknown;
  post: {
    id: string;
    loginID: string;
    title: string;
    content: string;
    createdAt: Date;
    updatedAt: Date;
    deletedAt: Date;
  };
}

export const FeaturedPost: React.FC<FeaturedPostProps> = (props) => {
  const { post } = props;

  const handleContainerClick = (event: React.MouseEvent<HTMLDivElement>) => {
    props.handleValueChange(post.id);
  };

  return (
    <Grid item xs={12} md={6}>
      <CardActionArea component="a">
        <Card sx={{ display: 'flex' }}>
          <CardContent sx={{ flex: 1 }} onClick={handleContainerClick}>
            <Typography component="h2" variant="h5">
              {post.title}
            </Typography>
            {/* <Typography variant="subtitle1" color="text.secondary">
              {post.createdAt}
            </Typography> */}
            <Typography variant="subtitle1" paragraph>
              {post.title}
            </Typography>
            <Typography variant="subtitle1" color="primary">
              Continue reading...
            </Typography>
          </CardContent>
          {/* <CardMedia
            component="img"
            sx={{ width: 160, display: { xs: 'none', sm: 'block' } }}
            image={post.image}
            alt={post.imageLabel}
          /> */}
        </Card>
      </CardActionArea>
    </Grid>
  );
}

// export default FeaturedPost;