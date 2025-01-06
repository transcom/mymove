/* eslint-disable camelcase */
import React from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as Yup from 'yup';

import AllowancesDetailForm from '../../../components/Office/AllowancesDetailForm/AllowancesDetailForm';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import { milmoveLogger } from 'utils/milmoveLog';
import { ORDERS_BRANCH_OPTIONS } from 'constants/orders';
import { ORDERS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { counselingUpdateAllowance } from 'services/ghcApi';
import { dropdownInputOptions } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const branchDropdownOption = dropdownInputOptions(ORDERS_BRANCH_OPTIONS);

const validationSchema = Yup.object({
  proGearWeight: Yup.number()
    .min(0, 'Pro-gear weight must be greater than or equal to 0')
    .max(2000, "Enter a weight that does not go over the customer's maximum allowance")
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  proGearWeightSpouse: Yup.number()
    .min(0, 'Spouse pro-gear weight must be greater than or equal to 0')
    .max(500, "Enter a weight that does not go over the customer's maximum allowance")
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  requiredMedicalEquipmentWeight: Yup.number()
    .min(0, 'Required medical equipment weight must be greater than or equal to 0')
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
  storageInTransit: Yup.number()
    .min(0, 'Storage in transit (days) must be greater than or equal to 0')
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
});

const ServicesCounselingMoveAllowances = () => {
  const { moveCode } = useParams();
  const navigate = useNavigate();

  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const orderId = move?.ordersId;

  const handleClose = () => {
    navigate(`../${servicesCounselingRoutes.MOVE_VIEW_PATH}`);
  };
  const queryClient = useQueryClient();
  const { mutate: mutateOrders } = useMutation(counselingUpdateAllowance, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.orders[variables.orderID];
      queryClient.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryClient.invalidateQueries(ORDERS);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const order = Object.values(orders)?.[0];
  const onSubmit = (values) => {
    const {
      grade,
      agency,
      dependentsAuthorized,
      proGearWeight,
      proGearWeightSpouse,
      requiredMedicalEquipmentWeight,
      organizationalClothingAndIndividualEquipment,
      storageInTransit,
      gunSafe,
      accompaniedTour,
      dependentsTwelveAndOver,
      dependentsUnderTwelve,
    } = values;
    const body = {
      issueDate: order.date_issued,
      newDutyLocationId: order.destinationDutyLocation.id,
      ordersNumber: order.order_number,
      ordersType: order.order_type,
      originDutyLocationId: order.originDutyLocation.id,
      reportByDate: order.report_by_date,
      grade,
      agency,
      dependentsAuthorized,
      proGearWeight: Number(proGearWeight),
      proGearWeightSpouse: Number(proGearWeightSpouse),
      requiredMedicalEquipmentWeight: Number(requiredMedicalEquipmentWeight),
      storageInTransit: Number(storageInTransit),
      organizationalClothingAndIndividualEquipment,
      gunSafe,
      accompaniedTour,
      dependentsTwelveAndOver: Number(dependentsTwelveAndOver),
      dependentsUnderTwelve: Number(dependentsUnderTwelve),
    };
    return mutateOrders({ orderID: orderId, ifMatchETag: order.eTag, body });
  };

  const { entitlement, grade, agency } = order;
  const {
    dependentsAuthorized,
    proGearWeight,
    proGearWeightSpouse,
    requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment,
    gunSafe,
    storageInTransit,
    dependentsUnderTwelve,
    dependentsTwelveAndOver,
    accompaniedTour,
  } = entitlement;

  const initialValues = {
    grade,
    agency,
    dependentsAuthorized,
    proGearWeight: `${proGearWeight}`,
    proGearWeightSpouse: `${proGearWeightSpouse}`,
    requiredMedicalEquipmentWeight: `${requiredMedicalEquipmentWeight}`,
    storageInTransit: `${storageInTransit}`,
    gunSafe,
    organizationalClothingAndIndividualEquipment,
    accompaniedTour,
    dependentsUnderTwelve: `${dependentsUnderTwelve}`,
    dependentsTwelveAndOver: `${dependentsTwelveAndOver}`,
  };

  return (
    <div className={styles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {(formik) => (
          <form onSubmit={formik.handleSubmit}>
            <div className={styles.content}>
              <div className={styles.top}>
                <Button
                  className={styles.closeButton}
                  data-testid="closeSidebar"
                  type="button"
                  onClick={handleClose}
                  unstyled
                >
                  <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
                </Button>
                <h2 className={styles.header} data-testid="allowances-header">
                  View allowances
                </h2>
                <div>
                  <Link className={styles.viewAllowances} data-testid="view-orders" to="../orders">
                    View orders
                  </Link>
                </div>
              </div>
              <div className={styles.body}>
                <AllowancesDetailForm
                  entitlements={order.entitlement}
                  branchOptions={branchDropdownOption}
                  header="Counseling"
                />
              </div>
              <div className={styles.bottom}>
                <div className={styles.buttonGroup}>
                  <Button disabled={formik.isSubmitting} data-testid="scAllowancesSave" type="submit">
                    Save
                  </Button>
                  <Button type="button" secondary onClick={handleClose}>
                    Cancel
                  </Button>
                </div>
              </div>
            </div>
          </form>
        )}
      </Formik>
    </div>
  );
};

export default ServicesCounselingMoveAllowances;
