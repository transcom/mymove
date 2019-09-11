import restProvider from './rest_provider';
import { fetchUtils, Admin, Resource, Layout } from 'react-admin';
import { createBrowserHistory } from 'history';
import React from 'react';
import Menu from './Menu';
import UserList from './UserList';
import UserCreate from './UserCreate';
import OfficeList from './OfficeList';
import ElectronicOrderList from './ElectronicOrderList';
import UserShow from './UserShow';
import styles from './Home.module.scss';
import { withContext } from 'shared/AppContext';

const httpClient = (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: 'application/json' });
  }
  // send cookies in the request
  options.credentials = 'same-origin';
  return fetchUtils.fetchJson(url, options);
};

const dataProvider = restProvider('/admin/v1', httpClient);
const AdminLayout = props => <Layout {...props} menu={Menu} />;
const history = createBrowserHistory({ basename: '/system' });

const Home = props => (
  <div className={styles['admin-system-wrapper']}>
    <Admin dataProvider={dataProvider} history={history} appLayout={AdminLayout}>
      <Resource
        name="office_users"
        options={{ label: 'Office users' }}
        list={UserList}
        show={UserShow}
        create={props.context.flags.createAdminUser && UserCreate}
      />
      <Resource name="offices" options={{ label: 'Offices' }} list={OfficeList} />
      <Resource name="electronic_orders" options={{ label: 'Electronic orders' }} list={ElectronicOrderList} />
    </Admin>
  </div>
);

const homeWithContext = withContext(Home);
export default homeWithContext;
