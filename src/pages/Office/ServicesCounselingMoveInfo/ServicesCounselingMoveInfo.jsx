import React, { Suspense, lazy, useState, useEffect } from 'react';
import { Switch, useParams, Redirect, Route, useHistory } from 'react-router-dom';

import 'styles/office.scss';
import CustomerHeader from 'components/CustomerHeader';
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
const CustomerInfo = lazy(() => import('pages/Office/CustomerInfo/CustomerInfo'));
const ServicesCounselingEditShipmentDetails = lazy(() =>
  import('pages/Office/ServicesCounselingEditShipmentDetails/ServicesCounselingEditShipmentDetails'),
);
const ServicesCounselingAddShipment = lazy(() =>
  import('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment'),
);

const ServicesCounselingMoveInfo = () => {
  const [customerEditAlert, setCustomerEditAlert] = useState(null);

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

      <Suspense fallback={<LoadingPlaceholder />}>
        <Switch>
          {/* TODO - Routes not finalized, revisit */}
          <Route path={servicesCounselingRoutes.MOVE_VIEW_PATH} exact>
            <ServicesCounselingMoveDetails editAlert={customerEditAlert} />
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

          <Route key="servicesCounselingAddShipment" path={servicesCounselingRoutes.SHIPMENT_ADD_PATH}>
            <ServicesCounselingAddShipment onUpdate={onCustomerInfoUpdate} />
          </Route>

          <Route key="servicesCounselingEditShipmentDetailsRoute" path={servicesCounselingRoutes.SHIPMENT_EDIT_PATH}>
            <ServicesCounselingEditShipmentDetails />
          </Route>

          {/* TODO - clarify role/tab access */}
          <Redirect from={servicesCounselingRoutes.BASE_MOVE_PATH} to={servicesCounselingRoutes.MOVE_VIEW_PATH} />
        </Switch>
      </Suspense>
    </>
  );
};

export default ServicesCounselingMoveInfo;
