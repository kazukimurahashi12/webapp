import axios from 'axios';
import ToastNotification from '../toast/toastNotification';


const api = axios.create({
  baseURL: 'http://localhost:8080/', // サーバーのAPIのベースURLを設定
});

// レスポンスインターセプターを追加してエラーハンドリングを行う
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 400) {
      const errorMessage = error.response.data.message; // エラーメッセージはサーバー側のレスポンスによって異なる場合があります
      ToastNotification({ message: errorMessage }); // トーストを表示
    }
    return Promise.reject(error);
  }
);

export default api;
