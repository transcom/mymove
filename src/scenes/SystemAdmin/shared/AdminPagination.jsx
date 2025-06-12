import React from 'react';
import { Pagination, useListContext } from 'react-admin';

import styles from './AdminPagination.module.scss';

const AdminPagination = () => {
  const { isLoading, total } = useListContext();
  return !isLoading && total === 0 ? (
    <div className={styles['no-results']}>No results found.</div>
  ) : (
    <Pagination rowsPerPageOptions={[]} isLoading={isLoading} total={total} />
  );
};
export default AdminPagination;
