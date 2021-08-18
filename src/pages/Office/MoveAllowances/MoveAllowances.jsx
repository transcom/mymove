/* eslint-disable camelcase */
import React from 'react';
import { Link, useHistory, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { queryCache, useMutation } from 'react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as Yup from 'yup';

import ordersStyles from '../Orders/Orders.module.scss';
import AllowancesDetailForm from '../../../components/Office/AllowancesDetailForm/AllowancesDetailForm';

import { MILMOVE_LOG_LEVEL } from 'shared/constants';
import { milmoveLog } from 'shared/milmoveLog';
import { updateAllowance } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { dropdownInputOptions } from 'shared/formatters';
import { ORDERS } from 'constants/queryKeys';

const rankDropdownOptions = dropdownInputOptions(ORDERS_RANK_OPTIONS);

const branchDropdownOption = dropdownInputOptions(ORDERS_BRANCH_OPTIONS);

const validationSchema = Yup.object({
  authorizedWeight: Yup.number().min(1, 'Authorized weight must be greater than or equal to 1').required('Required'),
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
    .min(0, 'RME weight must be greater than or equal to 0')
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
});

const MoveAllowances = () => {
  const { moveCode } = useParams();
  const history = useHistory();

  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const orderId = move?.ordersId;

  const handleClose = () => {
    history.push(`/moves/${moveCode}/details`);
  };

  const [mutateOrders] = useMutation(updateAllowance, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.orders[variables.orderID];
      queryCache.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryCache.invalidateQueries([ORDERS, variables.orderID]);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const order = Object.values(orders)?.[0];
  const onSubmit = (values) => {
    const {
      grade,
      authorizedWeight,
      agency,
      dependentsAuthorized,
      proGearWeight,
      proGearWeightSpouse,
      requiredMedicalEquipmentWeight,
      organizationalClothingAndIndividualEquipment,
    } = values;
    const body = {
      issueDate: order.date_issued,
      newDutyStationId: order.destinationDutyStation.id,
      ordersNumber: order.order_number,
      ordersType: order.order_type,
      originDutyStationId: order.originDutyStation.id,
      reportByDate: order.report_by_date,
      grade,
      authorizedWeight: Number(authorizedWeight),
      agency,
      dependentsAuthorized,
      proGearWeight: Number(proGearWeight),
      proGearWeightSpouse: Number(proGearWeightSpouse),
      requiredMedicalEquipmentWeight: Number(requiredMedicalEquipmentWeight),
      organizationalClothingAndIndividualEquipment,
    };
    mutateOrders({ orderID: orderId, ifMatchETag: order.eTag, body });
  };

  const { entitlement, grade, agency } = order;
  const {
    authorizedWeight,
    dependentsAuthorized,
    proGearWeight,
    proGearWeightSpouse,
    requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment,
  } = entitlement;

  const initialValues = {
    authorizedWeight: `${authorizedWeight}`,
    grade,
    agency,
    dependentsAuthorized,
    proGearWeight: `${proGearWeight}`,
    proGearWeightSpouse: `${proGearWeightSpouse}`,
    requiredMedicalEquipmentWeight: `${requiredMedicalEquipmentWeight}`,
    organizationalClothingAndIndividualEquipment,
  };

  return (
    <div className={ordersStyles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {(formik) => (
          <form onSubmit={formik.handleSubmit}>
            <div className={ordersStyles.orderDetails}>
              <div className={ordersStyles.top}>
                <Button
                  className={ordersStyles.closeButton}
                  data-testid="closeSidebar"
                  type="button"
                  onClick={handleClose}
                  unstyled
                >
                  <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
                </Button>
                <h2 className={ordersStyles.header} data-testid="allowances-header">
                  View Allowances
                </h2>
                <div>
                  <Link className={ordersStyles.viewAllowances} data-testid="view-orders" to="orders">
                    View Orders
                  </Link>
                </div>
              </div>
              <div className={ordersStyles.body}>
                <AllowancesDetailForm
                  entitlements={order.entitlement}
                  rankOptions={rankDropdownOptions}
                  branchOptions={branchDropdownOption}
                  editableAuthorizedWeight
                />
              </div>
              <div className={ordersStyles.bottom}>
                <div className={ordersStyles.buttonGroup}>
                  <Button disabled={formik.isSubmitting} type="submit">
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

export default MoveAllowances;
