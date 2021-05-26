import React, { Suspense, lazy, useState, useEffect } from 'react';
import { Switch, useParams, Redirect, Route, useHistory } from 'react-router-dom';

import 'styles/office.scss';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import CustomerHeader from 'components/CustomerHeader';
import { useTXOMoveInfoQueries } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingMoveDocumentWrapper = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDocumentWrapper/ServicesCounselingMoveDocumentWrapper'),
);
const ServicesCounselingMoveDetails = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDetails/ServicesCounselingMoveDetails'),
);
const CustomerInfo = lazy(() => import('pages/Office/CustomerInfo/CustomerInfo'));

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
          <Route path="/counseling/moves/:moveCode/details" exact>
            <ServicesCounselingMoveDetails customerEditAlert={customerEditAlert} />
          </Route>

          <Route path={['/counseling/moves/:moveCode/allowances', '/counseling/moves/:moveCode/orders']} exact>
            <ServicesCounselingMoveDocumentWrapper />
          </Route>

          <Route path="/counseling/moves/:moveCode/customer" exact>
            <CustomerInfo
              ordersId={order.id}
              customer={customerData}
              isLoading={isLoading}
              isError={isError}
              onUpdate={onCustomerInfoUpdate}
            />
          </Route>

          {/* TODO - clarify role/tab access */}
          <Redirect from="/counseling/moves/:moveCode" to="/counseling/moves/:moveCode/details" />
        </Switch>
      </Suspense>
    </>
  );
};

export default ServicesCounselingMoveInfo;
