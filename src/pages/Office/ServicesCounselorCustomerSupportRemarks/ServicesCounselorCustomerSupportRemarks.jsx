import React from 'react';
import 'styles/office.scss';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import customerSupportRemarkStyles from './ServicesCounselorCustomerSupportRemarks.module.scss';

import { useCustomerSupportRemarksQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import CustomerSupportRemarkText from 'components/Office/CustomerSupportRemarkText/CustomerSupportRemarkText';

const ServicesCounselorCustomerSupportRemarks = () => {
  const { moveCode } = useParams();
  const { customerSupportRemarks, isLoading, isError } = useCustomerSupportRemarksQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={classnames(styles.tabContent, customerSupportRemarkStyles.customerSupportRemarksContent)}>
      <GridContainer className={customerSupportRemarkStyles.customerSupportRemarksTitle}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <h1>Customer Support Remarks</h1>
          </Grid>
        </Grid>
      </GridContainer>
      <GridContainer>
        <Grid row>
          <Grid
            className={customerSupportRemarkStyles.customerSupportRemarksContainer}
            col
            desktop={{ col: 8, offset: 2 }}
          >
            <h2>Remarks</h2>
            <h3>Past Remarks</h3>
            {customerSupportRemarks.length === 0 && <p>No remarks yet.</p>}
            {customerSupportRemarks.length > 0 &&
              customerSupportRemarks.map((customerSupportRemark) => {
                return (
                  <CustomerSupportRemarkText
                    customerSupportRemark={customerSupportRemark}
                    key={customerSupportRemark.id}
                  />
                );
              })}
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};
export default ServicesCounselorCustomerSupportRemarks;
