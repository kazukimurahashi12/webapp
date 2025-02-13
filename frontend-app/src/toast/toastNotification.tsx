import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const ToastNotification = ({ message }) => {
  return toast.error(message);
};

export default ToastNotification;
