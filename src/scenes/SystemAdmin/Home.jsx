import { Admin, AppBar, fetchUtils, Layout, Resource, CustomRoutes, defaultTheme, TitlePortal } from 'react-admin';
import { Route } from 'react-router-dom';
import React from 'react';
import Cookies from 'js-cookie';

import WebhookSubscriptionEdit from '../../pages/Admin/WebhookSubscriptions/WebhookSubscriptionEdit';

import AdminLogoutOnInactivity from 'layout/AdminIdleTimeout';

import restProvider from './shared/rest_provider';
import Menu from './shared/Menu';
import UploadShow from './Uploads/UploadShow';
import ClientCertList from 'pages/Admin/ClientCerts/ClientCertList';
import ClientCertShow from 'pages/Admin/ClientCerts/ClientCertShow';
import ClientCertCreate from 'pages/Admin/ClientCerts/ClientCertCreate';
import ClientCertEdit from 'pages/Admin/ClientCerts/ClientCertEdit';
import ElectronicOrderList from './ElectronicOrders/ElectronicOrderList';
import FeatureFlagEvaluate from 'pages/Admin/FeatureFlags/FeatureFlagEvaluate';
import styles from './Home.module.scss';
import NotificationList from './Notifications/NotificationList';
import UploadSearch from './Uploads/UploadSearch';

import { milmoveLogger } from 'utils/milmoveLog';
import OfficeUserList from 'pages/Admin/OfficeUsers/OfficeUserList';
import OfficeUserShow from 'pages/Admin/OfficeUsers/OfficeUserShow';
import OfficeUserCreate from 'pages/Admin/OfficeUsers/OfficeUserCreate';
import OfficeUserEdit from 'pages/Admin/OfficeUsers/OfficeUserEdit';
import AdminUserList from 'pages/Admin/AdminUsers/AdminUserList';
import AdminUserShow from 'pages/Admin/AdminUsers/AdminUserShow';
import AdminUserCreate from 'pages/Admin/AdminUsers/AdminUserCreate';
import AdminUserEdit from 'pages/Admin/AdminUsers/AdminUserEdit';
import OfficeList from 'pages/Admin/Offices/OfficeList';
import MoveList from 'pages/Admin/Moves/MoveList';
import MoveShow from 'pages/Admin/Moves/MoveShow';
import MoveEdit from 'pages/Admin/Moves/MoveEdit';
import UserList from 'pages/Admin/Users/UserList';
import UserShow from 'pages/Admin/Users/UserShow';
import UserEdit from 'pages/Admin/Users/UserEdit';
import WebhookSubscriptionList from 'pages/Admin/WebhookSubscriptions/WebhookSubscriptionsList';
import WebhookSubscriptionShow from 'pages/Admin/WebhookSubscriptions/WebhookSubscriptionShow';
import WebhookSubscriptionCreate from 'pages/Admin/WebhookSubscriptions/WebhookSubscriptionCreate';
import RequestedOfficeUserList from 'pages/Admin/RequestedOfficeUsers/RequestedOfficeUserList';
import RequestedOfficeUserShow from 'pages/Admin/RequestedOfficeUsers/RequestedOfficeUserShow';
import RequestedOfficeUserEdit from 'pages/Admin/RequestedOfficeUsers/RequestedOfficeUserEdit';
import RejectedOfficeUserList from 'pages/Admin/RejectedOfficeUsers/RejectedOfficeUserList';
import RejectedOfficeUserShow from 'pages/Admin/RejectedOfficeUsers/RejectedOfficeUserShow';
import PaymentRequest858List from 'pages/Admin/PaymentRequests/PaymentRequest858List';
import PaymentRequest858Show from 'pages/Admin/PaymentRequests/PaymentRequest858Show';
import EdiErrorsList from 'pages/Admin/PaymentRequests/EdiErrorsList';
import EdiErrorsShow from 'pages/Admin/PaymentRequests/EdiErrorsShow';
import { deepmerge } from '@mui/utils';

import logo from '../../shared/images/milmove-logo.svg';
import { Box } from '@material-ui/core';

const httpClient = (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: 'application/json' });
  }
  const token = Cookies.get('masked_gorilla_csrf');
  if (!token) {
    milmoveLogger.warn('Unable to retrieve CSRF Token from cookie');
  }
  options.headers.set('X-CSRF-TOKEN', token);
  // send cookies in the request
  options.credentials = 'same-origin';
  return fetchUtils.fetchJson(url, options);
};

const CustomAppBar = (props) => (
  <AppBar {...props} position="sticky">
    {/* Logo on the left */}
    <Box component="span" sx={{ display: 'flex', alignItems: 'center', pl: 2 }}>
      <img src={logo} alt="logo" style={{ height: 22, paddingRight: '25px' }} />
    </Box>
    <TitlePortal />
  </AppBar>
);

const CUIWrapper = () => (
  <>
    <AdminLogoutOnInactivity />
    <CustomAppBar />
  </>
);

const dataProvider = restProvider('/admin/v1', httpClient);
const AdminLayout = (props) => <Layout {...props} menu={Menu} appBar={CUIWrapper} />;

const theme = deepmerge(defaultTheme, {
  palette: {
    secondary: {
      main: '#252f3e',
    },
  },
  components: {
    RaButton: {
      styleOverrides: {
        root: {
          backgroundColor: '#005ea2',
          height: '35px',
          padding: '8px',
          color: '#FFF !important',
          '&:hover': {
            backgroundColor: '#005ea2',
            opacity: '0.8',
          },
          '&:visited': {
            color: '##2491ff !important',
          },
        },
      },
    },
    RaMenuItemLink: {
      styleOverrides: {
        root: {
          color: '#000 !important',
          '&.RaMenuItemLink-active': {
            color: '#000 !important',
            backgroundColor: 'transparent !important',
          },
          '&:hover': {
            color: '#000 !important',
          },
        },
      },
    },
    RaSaveButton: {
      styleOverrides: {
        root: {
          minHeight: '30px',
          backgroundColor: '#005ea2',
          color: '#fff',
        },
      },
    },
    RaConfirm: {
      styleOverrides: {
        root: {
          '& .RaConfirm-confirmPrimary': {
            minHeight: '30px',
            backgroundColor: '#005ea2',
            color: '#fff',
            '&:hover': {
              opacity: '0.8',
              backgroundColor: '#1a4480',
            },
          },
          '& .MuiDialogActions-root > button:nth-of-type(1)': {
            minHeight: '30px',
            color: '#005ea2',
            border: '1px solid #005ea2',
            backgroundColor: 'transparent',
            '&:hover': {
              backgroundColor: '#f5f5f5',
              borderColor: '#004a8c',
            },
          },
        },
      },
    },
  },
});

const Home = () => (
  <div className={styles['admin-system-wrapper']}>
    <Admin dataProvider={dataProvider} basename="/system" layout={AdminLayout} disableTelemetry theme={theme}>
      <Resource
        name="requested-office-users"
        options={{ label: 'Requested Office Users' }}
        list={RequestedOfficeUserList}
        show={RequestedOfficeUserShow}
        edit={RequestedOfficeUserEdit}
      />
      <Resource
        name="rejected-office-users"
        options={{ label: 'Rejected Office Users' }}
        list={RejectedOfficeUserList}
        show={RejectedOfficeUserShow}
      />
      <Resource
        name="office-users"
        options={{ label: 'Office Users' }}
        list={OfficeUserList}
        show={OfficeUserShow}
        create={OfficeUserCreate}
        edit={OfficeUserEdit}
      />
      <Resource name="offices" options={{ label: 'Offices' }} list={OfficeList} />
      <Resource
        name="admin-users"
        options={{ label: 'Admin Users' }}
        list={AdminUserList}
        show={AdminUserShow}
        create={AdminUserCreate}
        edit={AdminUserEdit}
      />
      <Resource name="users" options={{ label: 'Users' }} list={UserList} show={UserShow} edit={UserEdit} />
      <Resource name="moves" options={{ label: 'Moves' }} list={MoveList} show={MoveShow} edit={MoveEdit} />
      <Resource
        name="payment-request-syncada-files"
        options={{ label: 'Payment Request Syncada Files' }}
        list={PaymentRequest858List}
        show={PaymentRequest858Show}
      />
      <Resource
        name="edi-errors"
        options={{ label: 'EDI Error Payment Requests' }}
        list={EdiErrorsList}
        show={EdiErrorsShow}
      />

      <Resource name="electronic-orders" options={{ label: 'Electronic orders' }} list={ElectronicOrderList} />
      <Resource name="uploads" options={{ label: 'Search Upload by ID' }} show={UploadShow} />
      <Resource name="organizations" />
      <Resource
        name="client-certificates"
        options={{ label: 'Client Certs' }}
        list={ClientCertList}
        show={ClientCertShow}
        create={ClientCertCreate}
        edit={ClientCertEdit}
      />
      <Resource name="notifications" options={{ label: 'Notifications' }} list={NotificationList} />
      <Resource name="feature-flags" options={{ label: 'Evaluate Feature Flag' }} list={FeatureFlagEvaluate} />
      <Resource
        name="webhook-subscriptions"
        options={{ label: 'Webhook Subscriptions' }}
        show={WebhookSubscriptionShow}
        create={WebhookSubscriptionCreate}
        list={WebhookSubscriptionList}
        edit={WebhookSubscriptionEdit}
      />
      <CustomRoutes>
        {/* Custom route for search by id for uploads */}
        <Route end path="/uploads" element={<UploadSearch />} />
        <Route end path="/feature-flags" element={<FeatureFlagEvaluate />} />
      </CustomRoutes>
    </Admin>
  </div>
);

export default Home;
