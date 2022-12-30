import React from 'react';
import { Pagination, useListContext } from 'react-admin';
import styles from './AdminPagination.module.scss';

const AdminPagination = () => {
  const { isLoading, total } = useListContext();
  return !isLoading && total === 0 ? (
    <div className={styles['no-results']}>
      There are no results for this access code. Please check your entry to make sure you entered the correct letter
      combination.
    </div>
  ) : (
    <Pagination rowsPerPageOptions={[]} isLoading={isLoading} total={total} />
  );
};
export default AdminPagination;
