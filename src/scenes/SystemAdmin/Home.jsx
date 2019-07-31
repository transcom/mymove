import restProvider from 'ra-data-simple-rest';
import { fetchUtils, Admin, Resource, Layout, List, Pagination, Datagrid, TextField } from 'react-admin';
import { history } from 'shared/store';
import React from 'react';
import Menu from './Menu';

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
const AdminPagination = props => <Pagination rowsPerPageOptions={[]} {...props} />;
const UserList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={500}>
    <Datagrid>
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="first_name" />
      <TextField source="last_name" />
    </Datagrid>
  </List>
);

const Home = () => (
  <div className="admin-system-wrapper">
    <Admin dataProvider={dataProvider} history={history} appLayout={AdminLayout}>
      <Resource name="office_users" list={UserList} />
    </Admin>
  </div>
);

export default Home;
