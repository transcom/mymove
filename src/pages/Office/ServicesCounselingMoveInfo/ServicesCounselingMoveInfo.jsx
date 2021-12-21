import React, { Suspense, lazy, useState, useEffect } from 'react';
import { Switch, useParams, Redirect, Route, useHistory } from 'react-router-dom';
import { useSelector } from 'react-redux';

import 'styles/office.scss';
import CustomerHeader from 'components/CustomerHeader';
import SystemError from 'components/SystemError';
import { servicesCounselingRoutes } from 'constants/routes';
import { useTXOMoveInfoQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingMoveDocumentWrapper = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDocumentWrapper/ServicesCounselingMoveDocumentWrapper'),
);
const ServicesCounselingMoveDetails = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDetails/ServicesCounselingMoveDetails'),
);
const ServicesCounselingEditShipmentDetails = lazy(() =>
  import('pages/Office/ServicesCounselingEditShipmentDetails/ServicesCounselingEditShipmentDetails'),
);
const CustomerInfo = lazy(() => import('pages/Office/CustomerInfo/CustomerInfo'));

const ServicesCounselingMoveInfo = () => {
  const [customerEditAlert, setCustomerEditAlert] = useState(null);
  const { hasRecentError, traceId } = useSelector((state) => state.interceptor);
  const onCustomerInfoUpdate = (alertType) => {
    if (alertType === 'error') {
      setCustomerEditAlert({
        alertType,
        message: 'Something went wrong, and your changes were not saved. Please try again later.',
      });
    } else {
      setCustomerEditAlert({
        alertType,
        message: 'Your changes were saved.',
      });
    }
  };

  const history = useHistory();
  useEffect(() => {
    // clear alert when route changes
    const unlisten = history.listen(() => {
      if (customerEditAlert) {
        setCustomerEditAlert(null);
      }
    });
    return () => {
      unlisten();
    };
  }, [history, customerEditAlert]);

  const { moveCode } = useParams();
  const { order, customerData, isLoading, isError } = useTXOMoveInfoQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <>
      <CustomerHeader order={order} customer={customerData} moveCode={moveCode} />
      {hasRecentError && (
        <SystemError>
          Something isn&apos;t working, but we&apos;re not sure what. Wait a minute and try again.
          <br />
          If that doesn&apos;t fix it, contact the{' '}
          <a href="https://move.mil/customer-service#technical-help-desk">Technical Help Desk</a> and give them this
          code: <strong>{traceId}</strong>
        </SystemError>
      )}
      <Suspense fallback={<LoadingPlaceholder />}>
        <Switch>
          {/* TODO - Routes not finalized, revisit */}
          <Route path={servicesCounselingRoutes.MOVE_VIEW_PATH} exact>
            <ServicesCounselingMoveDetails customerEditAlert={customerEditAlert} />
          </Route>

          <Route
            path={[servicesCounselingRoutes.ALLOWANCES_EDIT_PATH, servicesCounselingRoutes.ORDERS_EDIT_PATH]}
            exact
          >
            <ServicesCounselingMoveDocumentWrapper />
          </Route>

          <Route path={servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH} exact>
            <CustomerInfo
              ordersId={order.id}
              customer={customerData}
              isLoading={isLoading}
              isError={isError}
              onUpdate={onCustomerInfoUpdate}
            />
          </Route>

          <Route
            path={servicesCounselingRoutes.SHIPMENT_EDIT_PATH}
            exact
            // eslint-disable-next-line react/jsx-props-no-spreading
            render={(props) => <ServicesCounselingEditShipmentDetails {...props} onUpdate={onCustomerInfoUpdate} />}
          />

          {/* TODO - clarify role/tab access */}
          <Redirect from={servicesCounselingRoutes.BASE_MOVE_PATH} to={servicesCounselingRoutes.MOVE_VIEW_PATH} />
        </Switch>
      </Suspense>
    </>
  );
};

export default ServicesCounselingMoveInfo;
