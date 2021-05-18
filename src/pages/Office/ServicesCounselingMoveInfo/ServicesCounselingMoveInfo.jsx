import React, { Suspense, lazy } from 'react';
import { Switch, useParams, Redirect, Route } from 'react-router-dom';

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
            <ServicesCounselingMoveDetails />
          </Route>

          <Route path={['/counseling/moves/:moveCode/allowances', '/counseling/moves/:moveCode/orders']} exact>
            <ServicesCounselingMoveDocumentWrapper />
          </Route>

          <Route path="/counseling/moves/:moveCode/customer-info" exact>
            <CustomerInfo customer={customerData} isLoading={isLoading} isError={isError} />
          </Route>

          {/* TODO - clarify role/tab access */}
          <Redirect from="/counseling/moves/:moveCode" to="/counseling/moves/:moveCode/details" />
        </Switch>
      </Suspense>
    </>
  );
};

export default ServicesCounselingMoveInfo;
