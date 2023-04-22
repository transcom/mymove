import React from 'react';
import { Pagination } from 'react-admin';

import styles from './AdminPagination.module.scss';

const AdminPagination = (props) =>
  !props.isLoading && props.total === 0 ? (
    <div className={styles['no-results']}>
      There are no results for this access code. Please check your entry to make sure you entered the correct letter
      combination.
    </div>
  ) : (
    <Pagination rowsPerPageOptions={[]} {...props} />
  );

export default AdminPagination;
