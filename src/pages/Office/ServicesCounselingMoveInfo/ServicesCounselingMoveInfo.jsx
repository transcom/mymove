import React, { Suspense, lazy } from 'react';
import { Switch, useParams, Redirect, Route } from 'react-router-dom';

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
          <Route path={servicesCounselingRoutes.MOVE_VIEW_PATH} exact>
            <ServicesCounselingMoveDetails />
          </Route>

          <Route
            path={[servicesCounselingRoutes.ALLOWANCES_EDIT_PATH, servicesCounselingRoutes.ORDERS_EDIT_PATH]}
            exact
          >
            <ServicesCounselingMoveDocumentWrapper />
          </Route>

          <Route path={servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH} exact>
            <CustomerInfo ordersId={order.id} customer={customerData} isLoading={isLoading} isError={isError} />
          </Route>

          {/* TODO - clarify role/tab access */}
          <Redirect from={servicesCounselingRoutes.BASE_MOVE_PATH} to={servicesCounselingRoutes.MOVE_VIEW_PATH} />
        </Switch>
      </Suspense>
    </>
  );
};

export default ServicesCounselingMoveInfo;
