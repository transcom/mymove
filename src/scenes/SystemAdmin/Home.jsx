import restProvider from './shared/rest_provider';
import { fetchUtils, Admin, Resource, Layout } from 'react-admin';
import { createBrowserHistory } from 'history';
import React from 'react';
import Menu from './shared/Menu';
import AccessCodeList from './AccessCodes/AccessCodeList';
import UploadShow from './Uploads/UploadShow';
import OfficeUserList from './OfficeUsers/OfficeUserList';
import OfficeUserCreate from './OfficeUsers/OfficeUserCreate';
import OfficeUserEdit from './OfficeUsers/OfficeUserEdit';
import OfficeUserShow from './OfficeUsers/OfficeUserShow';
import AdminUserList from './AdminUsers/AdminUserList';
import AdminUserShow from './AdminUsers/AdminUserShow';
import AdminUserCreate from './AdminUsers/AdminUserCreate';
import OfficeList from './Offices/OfficeList';
import TSPPList from './TSPPs/TSPPList';
import TSPPShow from './TSPPs/TSPPShow';
import ElectronicOrderList from './ElectronicOrders/ElectronicOrderList';
import styles from './Home.module.scss';
import * as Cookies from 'js-cookie';
import customRoutes from './CustomRoutes';
import AdminUserEdit from './AdminUsers/AdminUserEdit';

const httpClient = (url, options = {}) => {
  const token = Cookies.get('masked_gorilla_csrf');
  if (!token) {
    console.warn('Unable to retrieve CSRF Token from cookie');
  }

  if (!options.headers) {
    options.headers = new Headers({ Accept: 'application/json', 'X-CSRF-TOKEN': token });
  }
  // send cookies in the request
  options.credentials = 'same-origin';
  return fetchUtils.fetchJson(url, options);
};

const dataProvider = restProvider('/admin/v1', httpClient);
const AdminLayout = props => <Layout {...props} menu={Menu} />;
const history = createBrowserHistory({ basename: '/system' });

const Home = () => (
  <div className={styles['admin-system-wrapper']}>
    <Admin dataProvider={dataProvider} history={history} appLayout={AdminLayout} customRoutes={customRoutes}>
      <Resource
        name="office_users"
        options={{ label: 'Office users' }}
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
    </Admin>
  </div>
);

export default Home;
