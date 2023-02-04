import React from 'react';
import { Route } from 'react-router-dom';
import UploadSearch from './Uploads/UploadSearch';
import { CustomRoutes } from 'react-admin';

const Routes = () => (
  <CustomRoutes>
    {/* Custom route for search by id for uploads */}
    <Route end path="/uploads" element={<UploadSearch />} />
  </CustomRoutes>
);

export default Routes;
