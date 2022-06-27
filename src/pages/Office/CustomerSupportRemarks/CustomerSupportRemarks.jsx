import React, { useState } from 'react';
import 'styles/office.scss';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';
import { queryCache, useMutation } from 'react-query';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import customerSupportRemarkStyles from './CustomerSupportRemarks.module.scss';

import ConnectedDeleteCustomerSupportRemarkConfirmationModal from 'components/ConfirmationModals/DeleteCustomerSupportRemarkConfirmationModal';
import { useCustomerSupportRemarksQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import CustomerSupportRemarkText from 'components/Office/CustomerSupportRemarkText/CustomerSupportRemarkText';
import CustomerSupportRemarkForm from 'components/Office/CustomerSupportRemarkForm/CustomerSupportRemarkForm';
import { CUSTOMER_SUPPORT_REMARKS } from 'constants/queryKeys';
import { deleteCustomerSupportRemark } from 'services/ghcApi';
import Alert from 'shared/Alert';

const CustomerSupportRemarks = () => {
  const { moveCode } = useParams();
  const [showDeletionSuccess, setShowDeletionSuccess] = useState(false);
  const [customerSupportRemarkIDToDelete, setCustomerSupportRemarkIDToDelete] = useState(null);
  const { customerSupportRemarks, isLoading, isError } = useCustomerSupportRemarksQueries(moveCode);
  const [deleteCustomerSupportRemarkMutation] = useMutation(deleteCustomerSupportRemark, {
    onSuccess: async () => {
      await queryCache.invalidateQueries([CUSTOMER_SUPPORT_REMARKS, moveCode]);
      setCustomerSupportRemarkIDToDelete(null);
      setShowDeletionSuccess(true);
    },
  });

  const onDelete = (remarkID) => {
    deleteCustomerSupportRemarkMutation({ customerSupportRemarkID: remarkID });
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <>
      <ConnectedDeleteCustomerSupportRemarkConfirmationModal
        isOpen={customerSupportRemarkIDToDelete !== null}
        customerSupportRemarkID={customerSupportRemarkIDToDelete}
        onClose={() => setCustomerSupportRemarkIDToDelete(null)}
        onSubmit={onDelete}
      />
      <div className={classnames(styles.tabContent, customerSupportRemarkStyles.customerSupportRemarksContent)}>
        <GridContainer className={customerSupportRemarkStyles.customerSupportRemarksTitle}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {showDeletionSuccess && <Alert type="success">Your remark has been deleted.</Alert>}
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
                      onDelete={setCustomerSupportRemarkIDToDelete}
                    />
                  );
                })}
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </>
  );
};
export default CustomerSupportRemarks;
