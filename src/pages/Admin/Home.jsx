import React from 'react';
import PropTypes from 'prop-types';
import { Admin, AppBar, Layout, Resource } from 'react-admin';
import { createBrowserHistory } from 'history';

import styles from './Home.module.scss';
import WebhookSubscriptionEdit from './WebhookSubscriptions/WebhookSubscriptionEdit';

import Menu from 'scenes/SystemAdmin/shared/Menu';
import FOUOHeader from 'components/FOUOHeader';
import AccessCodeList from 'scenes/SystemAdmin/AccessCodes/AccessCodeList';
import UploadShow from 'scenes/SystemAdmin/Uploads/UploadShow';
import OfficeUserList from 'pages/Admin/OfficeUsers/OfficeUserList';
import OfficeUserShow from 'pages/Admin/OfficeUsers/OfficeUserShow';
import OfficeUserCreate from 'pages/Admin/OfficeUsers/OfficeUserCreate';
import OfficeUserEdit from 'pages/Admin/OfficeUsers/OfficeUserEdit';
import AdminUserList from 'pages/Admin/AdminUsers/AdminUserList';
import AdminUserShow from 'pages/Admin/AdminUsers/AdminUserShow';
import AdminUserCreate from 'pages/Admin/AdminUsers/AdminUserCreate';
import AdminUserEdit from 'pages/Admin/AdminUsers/AdminUserEdit';
import OfficeList from 'pages/Admin/Offices/OfficeList';
import TSPPList from 'scenes/SystemAdmin/TSPPs/TSPPList';
import TSPPShow from 'scenes/SystemAdmin/TSPPs/TSPPShow';
import ElectronicOrderList from 'scenes/SystemAdmin/ElectronicOrders/ElectronicOrderList';
import MoveList from 'pages/Admin/Moves/MoveList';
import MoveShow from 'pages/Admin/Moves/MoveShow';
import MoveEdit from 'pages/Admin/Moves/MoveEdit';
import UserList from 'pages/Admin/Users/UserList';
import UserShow from 'pages/Admin/Users/UserShow';
import UserEdit from 'pages/Admin/Users/UserEdit';
import WebhookSubscriptionList from 'pages/Admin/WebhookSubscriptions/WebhookSubscriptionsList';
import WebhookSubscriptionShow from 'pages/Admin/WebhookSubscriptions/WebhookSubscriptionShow';
import WebhookSubscriptionCreate from 'pages/Admin/WebhookSubscriptions/WebhookSubscriptionCreate';
import customRoutes from 'scenes/SystemAdmin/CustomRoutes';
import NotificationList from 'scenes/SystemAdmin/Notifications/NotificationList';

const FOUOWrapper = () => (
  <>
    <FOUOHeader />
    <AppBar />
  </>
);

const AdminLayout = (props) => <Layout {...props} menu={Menu} appBar={FOUOWrapper} />;

const Home = ({ basename, dataProvider }) => {
  const history = createBrowserHistory({ basename });

  return (
    <div className={styles['admin-system-wrapper']}>
      <Admin
        dataProvider={dataProvider}
        history={history}
        layout={AdminLayout}
        customRoutes={customRoutes}
        disableTelemetry
      >
        <Resource
          name="office_users"
          options={{ label: 'Office Users' }}
          list={OfficeUserList}
          show={OfficeUserShow}
          create={OfficeUserCreate}
          edit={OfficeUserEdit}
        />
        <Resource name="offices" options={{ label: 'Offices' }} list={OfficeList} />
        <Resource
          name="admin_users"
          options={{ label: 'Admin Users' }}
          list={AdminUserList}
          show={AdminUserShow}
          create={AdminUserCreate}
          edit={AdminUserEdit}
        />
        <Resource name="users" options={{ label: 'Users' }} list={UserList} show={UserShow} edit={UserEdit} />
        <Resource name="moves" options={{ label: 'Moves' }} list={MoveList} show={MoveShow} edit={MoveEdit} />
        <Resource
          name="transportation_service_provider_performances"
          options={{ label: 'TSPPs' }}
          list={TSPPList}
          show={TSPPShow}
        />
        <Resource name="electronic_orders" options={{ label: 'Electronic orders' }} list={ElectronicOrderList} />
        <Resource name="access_codes" options={{ label: 'Access codes' }} list={AccessCodeList} />
        <Resource name="uploads" options={{ label: 'Search Upload by ID' }} show={UploadShow} />
        <Resource name="organizations" />
        <Resource name="notifications" options={{ label: 'Notifications' }} list={NotificationList} />
        <Resource
          name="webhook_subscriptions"
          options={{ label: 'Webhook Subscriptions' }}
          show={WebhookSubscriptionShow}
          create={WebhookSubscriptionCreate}
          list={WebhookSubscriptionList}
          edit={WebhookSubscriptionEdit}
        />
      </Admin>
    </div>
  );
};

Home.propTypes = {
  basename: PropTypes.string.isRequired,
  dataProvider: PropTypes.oneOfType([PropTypes.func, PropTypes.shape({})]).isRequired,
};

export default Home;
