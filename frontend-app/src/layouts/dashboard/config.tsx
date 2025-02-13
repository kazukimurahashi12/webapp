import ChartBarIcon from '@heroicons/react/24/solid/ChartBarIcon';
import DocumentIcon from '@heroicons/react/24/solid/DocumentIcon';
import CogIcon from '@heroicons/react/24/solid/CogIcon';
import LockClosedIcon from '@heroicons/react/24/solid/LockClosedIcon';
import ShoppingBagIcon from '@heroicons/react/24/solid/ShoppingBagIcon';
import UserIcon from '@heroicons/react/24/solid/UserIcon';
import UserPlusIcon from '@heroicons/react/24/solid/UserPlusIcon';
import UsersIcon from '@heroicons/react/24/solid/UsersIcon';
import XCircleIcon from '@heroicons/react/24/solid/XCircleIcon';
import { SvgIcon } from '@mui/material';

export const itemsLogin = [
  // {
  //   title: 'Overview',
  //   path: '/',
  //   icon: (
  //     <SvgIcon fontSize="small">
  //       <ChartBarIcon />
  //     </SvgIcon>
  //   )
  // },
  {
    title: 'Blog',
    path: '/auth/overview',
    icon: (
      <SvgIcon fontSize="small">
        <ChartBarIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Account',
    path: '/auth/account',
    icon: (
      <SvgIcon fontSize="small">
        <UserIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Settings',
    path: '/auth/settings',
    icon: (
      <SvgIcon fontSize="small">
        <CogIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Register',
    path: '/auth/register',
    icon: (
      <SvgIcon fontSize="small">
        <UserPlusIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Logout',
    path: '/auth/logout',
    icon: (
      <SvgIcon fontSize="small">
        <LockClosedIcon />
      </SvgIcon>
    )
  }
];

export const itemsLogout = [
  // {
  //   title: 'Overview',
  //   path: '/',
  //   icon: (
  //     <SvgIcon fontSize="small">
  //       <ChartBarIcon />
  //     </SvgIcon>
  //   )
  // },
  {
    title: 'Blog',
    path: '/',
    icon: (
      <SvgIcon fontSize="small">
        <DocumentIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Account',
    path: '/auth/account',
    icon: (
      <SvgIcon fontSize="small">
        <UserIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Settings',
    path: '/auth/settings',
    icon: (
      <SvgIcon fontSize="small">
        <CogIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Register',
    path: '/auth/register',
    icon: (
      <SvgIcon fontSize="small">
        <UserPlusIcon />
      </SvgIcon>
    )
  },
  {
    title: 'Login',
    path: '/auth/login',
    icon: (
      <SvgIcon fontSize="small">
        <LockClosedIcon />
      </SvgIcon>
    )
  }
];