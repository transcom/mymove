import restProvider from './shared/rest_provider';
import { fetchUtils, Admin, Resource, Layout } from 'react-admin';
import { createBrowserHistory } from 'history';
import React from 'react';
import Menu from './shared/Menu';
import AccessCodeList from './AccessCodes/AccessCodeList';
import UserList from './OfficeUsers/UserList';
import UserCreate from './OfficeUsers/UserCreate';
import UserEdit from './OfficeUsers/UserEdit';
import UserShow from './OfficeUsers/UserShow';
import OfficeList from './Offices/OfficeList';
import ElectronicOrderList from './ElectronicOrders/ElectronicOrderList';
import styles from './Home.module.scss';
import * as Cookies from 'js-cookie';

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
    <Admin dataProvider={dataProvider} history={history} appLayout={AdminLayout}>
      <Resource
        name="office_users"
        options={{ label: 'Office users' }}
        list={UserList}
        show={UserShow}
        create={UserCreate}
        edit={UserEdit}
      />
      <Resource name="offices" options={{ label: 'Offices' }} list={OfficeList} />
      <Resource name="electronic_orders" options={{ label: 'Electronic orders' }} list={ElectronicOrderList} />
      <Resource name="access_codes" options={{ label: 'Access codes' }} list={AccessCodeList} />
    </Admin>
  </div>
);

export default Home;
