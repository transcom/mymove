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
import CustomerSupportRemarkForm from 'components/Office/CustomerSupportRemarkForm/CustomerSupportRemarkForm';

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
            <h1>Customer support remarks</h1>
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

            <CustomerSupportRemarkForm />

            <h3>Past remarks</h3>

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
