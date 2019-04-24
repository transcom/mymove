import React from 'react';
import { Provider } from 'react-redux';
import { Admin, Resource } from 'react-admin';
import defaultMessages from 'ra-language-english';
import restProvider from 'ra-data-simple-rest';
import adminReducer from './adminReducer';
import { history } from 'shared/store';
import customRoutes from './customRoutes';

const dataProvider = restProvider('http://admin/v1/...');
const i18nProvider = () => defaultMessages;

const AdminWrapper = () => (
  <Provider store={adminReducer({ dataProvider, i18nProvider, history })}>
    <Admin customRoutes={customRoutes} dataProvider={dataProvider} history={history}>
      <Resource />
    </Admin>
  </Provider>
);

export default AdminWrapper;
