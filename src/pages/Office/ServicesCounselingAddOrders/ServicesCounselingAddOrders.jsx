import React from 'react';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import styles from './ServicesCounselingAddOrders.module.scss';

import { dropdownInputOptions, formatYesNoAPIValue } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import AddOrdersForm from 'components/Office/AddOrdersForm/AddOrdersForm';
import { counselingCreateOrder } from 'services/ghcApi';
import { ORDERS } from 'constants/queryKeys';
import { formatDateForSwagger } from 'shared/dates';
import { servicesCounselingRoutes } from 'constants/routes';
import { milmoveLogger } from 'utils/milmoveLog';

const ServicesCounselingAddOrders = () => {
  const { customerId } = useParams();
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

  const ordersTypeOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
  const initialValues = {
    ordersType: '',
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
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export default ServicesCounselingAddOrders;
