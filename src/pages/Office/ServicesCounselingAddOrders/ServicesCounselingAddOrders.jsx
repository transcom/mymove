import React, { useEffect, useState } from 'react';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { generatePath, useLocation, useNavigate, useParams } from 'react-router-dom';

import styles from './ServicesCounselingAddOrders.module.scss';

import { dropdownInputOptions, formatYesNoAPIValue } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import AddOrdersForm from 'components/Office/AddOrdersForm/AddOrdersForm';
import { counselingCreateOrder } from 'services/ghcApi';
import { ORDERS } from 'constants/queryKeys';
import { formatDateForSwagger } from 'shared/dates';
import { servicesCounselingRoutes } from 'constants/routes';
import { milmoveLogger } from 'utils/milmoveLog';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

const ServicesCounselingAddOrders = ({ userPrivileges }) => {
  const { customerId } = useParams();
  const { state } = useLocation();
  const isSafetyMoveSelected = state?.isSafetyMoveSelected;
  const navigate = useNavigate();
  const handleBack = () => {
    navigate(-1);
  };
  const handleClose = (moveCode) => {
    const path = generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, {
      moveCode,
    });
    navigate(path);
  };
  const queryClient = useQueryClient();
  const { mutate: mutateOrders } = useMutation(counselingCreateOrder, {
    onSuccess: (data) => {
      const orderID = Object.keys(data.orders)[0];
      const updatedOrder = data.orders[orderID];
      queryClient.setQueryData([ORDERS, orderID], {
        orders: {
          [`${orderID}`]: updatedOrder,
        },
      });

      queryClient.invalidateQueries(ORDERS);
      handleClose(updatedOrder.moveCode);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const [isSafetyMoveFF, setSafetyMoveFF] = useState(false);

  useEffect(() => {
    isBooleanFlagEnabled('safety_move').then((enabled) => {
      setSafetyMoveFF(enabled);
    });
  }, []);

  const isSafetyPrivileged = isSafetyMoveFF
    ? userPrivileges?.some((privilege) => privilege.privilegeType === elevatedPrivilegeTypes.SAFETY)
    : false;

  const allowedOrdersTypes = isSafetyPrivileged
    ? { ...ORDERS_TYPE_OPTIONS, ...{ SAFETY: 'Safety' } }
    : ORDERS_TYPE_OPTIONS;
  const ordersTypeOptions = dropdownInputOptions(allowedOrdersTypes);

  const initialValues = {
    ordersType: isSafetyMoveSelected ? 'SAFETY' : '',
    issueDate: '',
    reportByDate: '',
    hasDependents: '',
    newDutyLocation: '',
    grade: '',
    originDutyLocation: '',
  };

  const handleSubmit = (values) => {
    const body = {
      ...values,
      serviceMemberId: customerId,
      newDutyLocationId: values.newDutyLocation.id,
      hasDependents: formatYesNoAPIValue(values.hasDependents),
      reportByDate: formatDateForSwagger(values.reportByDate),
      issueDate: formatDateForSwagger(values.issueDate),
      grade: values.grade,
      originDutyLocationId: values.originDutyLocation.id,
      spouseHasProGear: false,
    };
    mutateOrders({ body });
  };

  return (
    <GridContainer data-testid="main-container">
      <Grid row className={styles.ordersFormContainer} data-testid="orders-form-container">
        <Grid col>
          <AddOrdersForm
            onSubmit={handleSubmit}
            ordersTypeOptions={ordersTypeOptions}
            initialValues={initialValues}
            onBack={handleBack}
            isSafetyMoveSelected={isSafetyMoveSelected}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export default ServicesCounselingAddOrders;
