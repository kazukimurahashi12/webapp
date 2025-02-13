import { useCallback, useEffect, useState } from 'react';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Divider,
  Stack,
  TextField
} from '@mui/material';
import axios from 'axios';
import { useRouter } from 'next/router';

type User = {
  changePassword: string;
  nowPassword: string;
};

export const SettingsPassword = () => {
  const router = useRouter();
  const [propsUser, setUserProps] = useState<User>({
    changePassword: "",
    nowPassword: ""
  });
  const [isPosting, setPosting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string>("");

  const [values, setValues] = useState({
    password: '',
    confirm: ''
  });

  // const handleChange = useCallback(
  //   (event) => {
  //     setValues((prevState) => ({
  //       ...prevState,
  //       [event.target.name]: event.target.value
  //     }));
  //   },
  //   []
  // );

  // const handleSubmit = useCallback(
  //   (event) => {
  //     event.preventDefault();
  //   },
  //   []
  // );

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
          nowPassword: reqNowId,
          changePassword: reqNowId // 新しい値をidフィールドに代入
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
    if (propsUser.changePassword === propsUser.nowPassword) {
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
    <form onSubmit={handlePost}>
      <Card>
        <CardHeader
          subheader="Update password"
          title="Password"
        />
        <Divider />
        <CardContent>
          <Stack
            spacing={3}
            sx={{ maxWidth: 400 }}
          >
            <TextField
              fullWidth
              label="Password"
              name="password"
              onChange={handlePost}
              type="password"
              value={values.password}
            />
            <TextField
              fullWidth
              label="Password (Confirm)"
              name="confirm"
              onChange={handlePost}
              type="password"
              value={values.confirm}
            />
          </Stack>
        </CardContent>
        <Divider />
        <CardActions sx={{ justifyContent: 'flex-end' }}>
          <Button variant="contained">
            Update
          </Button>
        </CardActions>
      </Card>
    </form>
  );
};
