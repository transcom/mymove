import React from 'react';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import customerRemarkStyles from './CustomerSupportRemarks.module.scss';

import { useCustomerSupportRemarksQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { formatCustomerSupportRemarksDate } from 'utils/formatters';

const CustomerSupportRemarks = () => {
  const { moveCode } = useParams();
  const { customerSupportRemarks, isLoading, isError } = useCustomerSupportRemarksQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={styles.tabContent}>
      <GridContainer className={customerRemarkStyles.customerRemarksContainer}>
        <h1>Customer Support Remarks</h1>
      </GridContainer>
      <GridContainer className={classnames(customerRemarkStyles.remarksContainer, 'container--popout')}>
        <h2>Remarks</h2>
        <h4>Past Remarks</h4>
        {customerSupportRemarks.length === 0 && <p>No remarks yet.</p>}
        {customerSupportRemarks.length > 0 &&
          customerSupportRemarks.map((customerSupportRemark) => {
            return (
              <div key={customerSupportRemark.id}>
                <p className={customerRemarkStyles.customerRemarkBody}>
                  <small>
                    <strong>
                      {customerSupportRemark.officeUserFirstName} {customerSupportRemark.officeUserLastName}
                    </strong>{' '}
                    {formatCustomerSupportRemarksDate(customerSupportRemark.createdAt)}
                  </small>
                </p>
                {customerSupportRemark.content}
              </div>
            );
          })}
      </GridContainer>
    </div>
  );
};
export default CustomerSupportRemarks;
