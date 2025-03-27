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
      // エラーメッセージを取得
      const errorMessage = error.response.data.message; 
      // トーストを表示
      ToastNotification({ message: errorMessage }); 
    }
    return Promise.reject(error);
  }
);

export default api;
