/* eslint-disable camelcase */
import React from 'react';
import { generatePath } from 'react-router';
import { Link, useHistory, useParams } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import { queryCache, useMutation } from 'react-query';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import * as Yup from 'yup';

import documentWrapperStyles from '../ServicesCounselingMoveDocumentWrapper/ServicesCounselingMoveDocumentWrapper.module.scss';
import AllowancesDetailForm from '../../../components/Office/AllowancesDetailForm/AllowancesDetailForm';

import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { ORDERS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { updateAllowance } from 'services/ghcApi';
import { dropdownInputOptions } from 'shared/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const rankDropdownOptions = dropdownInputOptions(ORDERS_RANK_OPTIONS);

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
    .min(0, 'RME weight must be greater than or equal to 0')
    .transform((value) => (Number.isNaN(value) ? 0 : value))
    .notRequired(),
});

const ServicesCounselingMoveAllowances = () => {
  const { moveCode } = useParams();
  const history = useHistory();

  const { move, orders, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const orderId = move?.ordersId;

  const handleClose = () => {
    history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
  };

  const [mutateOrders] = useMutation(updateAllowance, {
    onSuccess: (data, variables) => {
      const updatedOrder = data.orders[variables.orderID];
      queryCache.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryCache.invalidateQueries(ORDERS);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message from unsuccessful mutation.
      // RA: TODO: As indicated, this error needs to be handled and needs further investigation and work.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5597
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      // eslint-disable-next-line no-console
      console.log(errorMsg);
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
    } = values;
    const body = {
      issueDate: order.date_issued,
      newDutyStationId: order.destinationDutyStation.id,
      ordersNumber: order.order_number,
      ordersType: order.order_type,
      originDutyStationId: order.originDutyStation.id,
      reportByDate: order.report_by_date,
      grade,
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
    dependentsAuthorized,
    proGearWeight,
    proGearWeightSpouse,
    requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment,
  } = entitlement;

  const initialValues = {
    grade,
    agency,
    dependentsAuthorized,
    proGearWeight: `${proGearWeight}`,
    proGearWeightSpouse: `${proGearWeightSpouse}`,
    requiredMedicalEquipmentWeight: `${requiredMedicalEquipmentWeight}`,
    organizationalClothingAndIndividualEquipment,
  };

  return (
    <div className={documentWrapperStyles.sidebar}>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {(formik) => (
          <form onSubmit={formik.handleSubmit}>
            <div className={documentWrapperStyles.orderDetails}>
              <div className={documentWrapperStyles.top}>
                <Button
                  className={documentWrapperStyles.closeButton}
                  data-testid="closeSidebar"
                  type="button"
                  onClick={handleClose}
                  unstyled
                >
                  <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
                </Button>
                <h2 className={documentWrapperStyles.header} data-testid="allowances-header">
                  View allowances
                </h2>
                <div>
                  <Link className={documentWrapperStyles.viewAllowances} data-testid="view-orders" to="orders">
                    View orders
                  </Link>
                </div>
              </div>
              <div className={documentWrapperStyles.body}>
                <AllowancesDetailForm
                  entitlements={order.entitlement}
                  rankOptions={rankDropdownOptions}
                  branchOptions={branchDropdownOption}
                  header="Counseling"
                />
              </div>
              <div className={documentWrapperStyles.bottom}>
                <div className={documentWrapperStyles.buttonGroup}>
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

export default ServicesCounselingMoveAllowances;
