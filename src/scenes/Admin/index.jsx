import React from 'react';
import { Admin, Resource } from 'react-admin';
import restProvider from 'ra-data-simple-rest';
import { history } from 'shared/store';
import customRoutes from './customRoutes';

const dataProvider = restProvider('http://admin/v1/...');

const AdminWrapper = () => (
  <div className="admin-system-wrapper">
    <Admin customRoutes={customRoutes} dataProvider={dataProvider} history={history}>
      <Resource />
    </Admin>
  </div>
);

export default AdminWrapper;
