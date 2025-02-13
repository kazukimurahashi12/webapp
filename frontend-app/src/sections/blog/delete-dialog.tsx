import React, { useState } from "react";
import { useRouter } from "next/router";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@mui/material";
import axios from "axios";

const DeleteDialog = ({ id }) => {
  const [open, setOpen] = useState(false);
  const [isDeleting, setDeleting] = useState(false);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleDelete = async () => {
    if (isDeleting) {
      //削除処理が実行中の場合は何もしない
      return;
    }
    setDeleting(true); //削除処理を実行中にセット

    const hostname = process.env.NODE_ENV === "production" ? "server-app" : "localhost";
    try {
      const response = await axios.get(`http://${hostname}:8080/blog/delete/${id}`, {
        withCredentials: true,
      });
      console.log("ブログ記事を削除しました");
      // 成功したことを示すステートを更新し、ダイアログを閉じる
      setOpen(false);
    } catch (error) {
      console.error("Error deleting blog:", error);
    } finally {
      // レンダリング後にリダイレクトするために非同期処理で setTimeout を使う
      setTimeout(() => {
        setDeleting(false);
        window.location.reload();
      }, 0);
    }
  };

  return (
    <>
      <Button onClick={handleOpen}>Delete</Button>
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>
          Are you sure you want to delete the blog post?(ブログ記事を削除しますか？)
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Deleting the blog post is irreversible. Are you sure you want to proceed with the
            deletion?(ブログ記事を削除すると、元に戻すことはできません。本当に削除しますか？)
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleDelete} variant="contained" color="error" disabled={isDeleting}>
            {isDeleting ? "Deleting..." : "Delete"} {/* 削除中はボタンを無効化 */}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default DeleteDialog;
