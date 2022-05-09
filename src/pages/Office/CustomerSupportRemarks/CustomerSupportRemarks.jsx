import React from 'react';
import 'styles/office.scss';
import { GridContainer } from '@trussworks/react-uswds';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import customerSupportRemarkStyles from './CustomerSupportRemarks.module.scss';

import { useCustomerSupportRemarksQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import CustomerSupportRemarkText from 'components/Office/CustomerSupportRemarkText/CustomerSupportRemarkText';

const CustomerSupportRemarks = () => {
  const { moveCode } = useParams();
  const { customerSupportRemarks, isLoading, isError } = useCustomerSupportRemarksQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={classnames(styles.tabContent, customerSupportRemarkStyles.customerSupportRemarksContent)}>
      <GridContainer className={customerSupportRemarkStyles.customerSupportRemarksTitle}>
        <h1>Customer Support Remarks</h1>
      </GridContainer>
      <GridContainer
        className={classnames(customerSupportRemarkStyles.customerSupportRemarksContainer, 'container--popout')}
      >
        <h2>Remarks</h2>
        <h4>Past Remarks</h4>
        {customerSupportRemarks.length === 0 && <p>No remarks yet.</p>}
        {customerSupportRemarks.length > 0 &&
          customerSupportRemarks.map((customerSupportRemark) => {
            return (
              <CustomerSupportRemarkText customerSupportRemark={customerSupportRemark} key={customerSupportRemark.id} />
            );
          })}
      </GridContainer>
    </div>
  );
};
export default CustomerSupportRemarks;
