import { useCallback } from 'react';
import { useRouter } from 'next/router';
import PropTypes from 'prop-types';
import { Box, Divider, MenuItem, MenuList, Popover, Typography } from '@mui/material';
import axios from 'axios';

import { useAuth } from '../../hooks/use-auth';
import React from 'react';


interface AccountPopoverProps {
  anchorEl: HTMLElement | null;
  onClose?: () => void;
  open: boolean;
}

export const AccountPopover: React.FC<AccountPopoverProps> = (props) => {
  const { anchorEl, onClose, open } = props;
  const router = useRouter();
  const auth = useAuth();

  // const handleSignOut = useCallback(() => {
  //   onClose?.();
  //   auth.signOut();
  //   router.push('/auth/login');
  // }, [onClose, auth, router]);

  const handleSignOut = (event: React.MouseEvent<HTMLElement>) => {
    const hostname = process.env.NODE_ENV === 'production' ? 'server-app' : 'localhost';

    axios
      .get(
        `http://${hostname}:8080/logout`,
        {
          // headers: {
          //   "Content-Type": "application/x-www-form-urlencoded",
          // },
          withCredentials: true,
        }
      )
      .then(() => {
        // ログアウト成功時の処理
        onClose?.();
        auth.signOut();
        router.push('/auth/login');
      })
      .catch((error) => {
        // ログイン失敗時の処理
        console.error(error);
      });
  };

  return (
    <Popover
      anchorEl={anchorEl}
      anchorOrigin={{
        horizontal: 'left',
        vertical: 'bottom'
      }}
      onClose={onClose}
      open={open}
      PaperProps={{ sx: { width: 200 } }}
    >
      <Box sx={{ py: 1.5, px: 2 }}>
        <Typography variant="overline">Account</Typography>
        <Typography color="text.secondary" variant="body2">
          Anika Visser
        </Typography>
      </Box>
      <Divider />
      <MenuList
        disablePadding
        dense
        sx={{
          p: '8px',
          '& > *': {
            borderRadius: 1
          }
        }}
      >
        <MenuItem onClick={handleSignOut}>Sign out</MenuItem>
      </MenuList>
    </Popover>
  );
};

AccountPopover.propTypes = {
  anchorEl: PropTypes.any,
  onClose: PropTypes.func,
  open: PropTypes.bool.isRequired
};
