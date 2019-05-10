import restProvider from 'ra-data-simple-rest';
import { Admin } from 'react-admin';
import { history } from 'shared/store';
import React from 'react';

const dataProvider = restProvider('http://admin/v1/...');

const AdminHome = () => (
  <div className="admin-system-wrapper">
    <Admin dataProvider={dataProvider} history={history}>
      {/*<Resource />*/}
    </Admin>
  </div>
);

export default AdminHome;
