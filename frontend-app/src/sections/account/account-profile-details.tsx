import { useCallback, useEffect, useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Divider,
  TextField,
  Unstable_Grid2 as Grid
} from '@mui/material';
import axios from 'axios';
import { useRouter } from "next/router";


type User = {
  currentId: string; //nowId
  newId: string; //changeId
}

export const AccountProfileDetails = () => {
  const router = useRouter();
  const [propsUser, setUserProps] = useState<User>({
    newId: "",
    currentId: ""
  });
  const [isPosting, setPosting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string>("");

  useEffect(() => {
    const fetchData = async () => {
      const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
      try {
        const response = await axios.get(`http://${hostname}:8080/api/login-id`, {
          withCredentials: true,
        });
        const reqNowId = response.data.id;
        setUserProps({
          ...propsUser, // 既存のuserIdオブジェクトを展開して新しいオブジェクトを作成
          currentId: reqNowId,
          newId: reqNowId // 新しい値をidフィールドに代入
        });
      } catch (error) {
        console.error(error);
      }
    };
    // コンポーネントのマウント時にリクエストを実行
    fetchData();
  }, []);

  const handlePost = async () => {
    //IDが既存のと同じ場合エラー
    if (propsUser.newId === propsUser.currentId) {
      setErrorMessage("Error: ID is the same as the current ID");
      return;
    }

    if (isPosting) {
      // 追加: 処理が実行中の場合は何もしない
      return;
    }
    setPosting(true); // 追加: 処理を実行中にセット

    const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
    try {
      const response = await axios.post(
        `http://${hostname}:8080/update/id`,
        propsUser, 
        {
          headers: {
            "Content-Type": "application/json", // JSON形式で送信するためのヘッダー設定application/json"でないと飛ばない
          },
          withCredentials: true,
        }
      );
      console.log("IDを変更しました");
    } catch (error) {
      // エラー処理を行う場合はここに記述
      console.error("Error editing id:", error);
    } finally {
      // 処理が終了したので false にセット
      setPosting(false);
      router.push("/auth/overview");
    }
  };

  return (
    <form
      autoComplete="off"
      noValidate
      onSubmit={handlePost}
    >
      <Card>
        <CardHeader
          subheader="ID can be edited"
          title="Profile"
        />
        <CardContent sx={{ pt: 0 }}>
          <Box sx={{ m: -1.5 }}>
            <Grid
              container
              spacing={3}
            >
            <Grid xs={12} md={6}>
              <TextField
                fullWidth
                label="ID"
                name="id"
                onChange={(event) => {
                  const inputId = event.target.value.trim();
                  setUserProps({
                    ...propsUser,
                    newId: inputId,
                  });
                  if (inputId === "") {
                    setErrorMessage("ID cannot be empty");
                  } else {
                    setErrorMessage("");
                  }
                }}
                required
                value={propsUser.newId}
                error={!!errorMessage}
                helperText={errorMessage}
              />
            </Grid>
              <Grid
                xs={12}
                md={6}
              >
              </Grid>
            </Grid>
          </Box>
        </CardContent>
        <Divider />
        <CardActions sx={{ justifyContent: 'flex-end' }}>
        <Button 
            variant="contained"
            onClick={handlePost}
            disabled={propsUser.newId === propsUser.currentId || !!errorMessage}
          >
            Save details
          </Button>
        </CardActions>
      </Card>
    </form>
  );
};
