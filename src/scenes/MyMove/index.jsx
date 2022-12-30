import React, { Component, lazy } from 'react';
import PropTypes from 'prop-types';
import { Route, Routes, Navigate } from 'react-router-dom';
import { connect } from 'react-redux';
import { GovBanner } from '@trussworks/react-uswds';

import 'styles/full_uswds.scss';
import 'styles/customer.scss';

import BypassBlock from 'components/BypassBlock';
import CUIHeader from 'components/CUIHeader/CUIHeader';
import LoggedOutHeader from 'containers/Headers/LoggedOutHeader';
import CustomerLoggedInHeader from 'containers/Headers/CustomerLoggedInHeader';
import Alert from 'shared/Alert';
import Footer from 'components/Customer/Footer';
import ConnectedLogoutOnInactivity from 'layout/LogoutOnInactivity';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { getWorkflowRoutes } from './getWorkflowRoutes';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import { withContext } from 'shared/AppContext';
import { no_op } from 'shared/utils';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { initOnboarding as initOnboardingAction } from 'store/onboarding/actions';
import { selectGetCurrentUserIsLoading, selectIsLoggedIn } from 'store/auth/selectors';
import { selectConusStatus } from 'store/onboarding/selectors';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentMove,
  selectHasCanceledMove,
} from 'store/entities/selectors';
import { generalRoutes, customerRoutes } from 'constants/routes';
/** Pages */
import InfectedUpload from 'shared/Uploader/InfectedUpload';
import ProcessingUpload from 'shared/Uploader/ProcessingUpload';
import PpmLanding from 'scenes/PpmLanding';
import Edit from 'scenes/Review/Edit';
import EditProfile from 'scenes/Review/EditProfile';
import EditDateAndLocation from 'scenes/Review/EditDateAndLocation';
import EditWeight from 'scenes/Review/EditWeight';
import PPMPaymentRequestIntro from 'scenes/Moves/Ppm/PPMPaymentRequestIntro';
import WeightTicket from 'scenes/Moves/Ppm/WeightTicket';
import ExpensesLanding from 'scenes/Moves/Ppm/ExpensesLanding';
import ExpensesUpload from 'scenes/Moves/Ppm/ExpensesUpload';
import AllowableExpenses from 'scenes/Moves/Ppm/AllowableExpenses';
import WeightTicketExamples from 'scenes/Moves/Ppm/WeightTicketExamples';
import NotFound from 'components/NotFound/NotFound';
import PrivacyPolicyStatement from 'shared/Statements/PrivacyAndPolicyStatement';
import AccessibilityStatement from 'shared/Statements/AccessibilityStatement';
import TrailerCriteria from 'scenes/Moves/Ppm/TrailerCriteria';
import PaymentReview from 'scenes/Moves/Ppm/PaymentReview/index';
import CustomerAgreementLegalese from 'scenes/Moves/Ppm/CustomerAgreementLegalese';
import ConnectedCreateOrEditMtoShipment from 'pages/MyMove/CreateOrEditMtoShipment';
import Home from 'pages/MyMove/Home';

// Pages should be lazy-loaded (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
const InvalidPermissions = lazy(() => import('pages/InvalidPermissions/InvalidPermissions'));
const MovingInfo = lazy(() => import('pages/MyMove/MovingInfo'));
const EditServiceInfo = lazy(() => import('pages/MyMove/Profile/EditServiceInfo'));
const Profile = lazy(() => import('pages/MyMove/Profile/Profile'));
const EditContactInfo = lazy(() => import('pages/MyMove/Profile/EditContactInfo'));
const AmendOrders = lazy(() => import('pages/MyMove/AmendOrders/AmendOrders'));
const EditOrders = lazy(() => import('pages/MyMove/EditOrders'));
const EstimatedWeightsProGear = lazy(() =>
  import('pages/MyMove/PPM/Booking/EstimatedWeightsProGear/EstimatedWeightsProGear'),
);
const EstimatedIncentive = lazy(() => import('pages/MyMove/PPM/Booking/EstimatedIncentive/EstimatedIncentive'));
const Advance = lazy(() => import('pages/MyMove/PPM/Booking/Advance/Advance'));
const About = lazy(() => import('pages/MyMove/PPM/Closeout/About/About'));
const WeightTickets = lazy(() => import('pages/MyMove/PPM/Closeout/WeightTickets/WeightTickets'));
const PPMReview = lazy(() => import('pages/MyMove/PPM/Closeout/Review/Review'));
const ProGear = lazy(() => import('pages/MyMove/PPM/Closeout/ProGear/ProGear.jsx'));
const Expenses = lazy(() => import('pages/MyMove/PPM/Closeout/Expenses/Expenses'));
const PPMFinalCloseout = lazy(() => import('pages/MyMove/PPM/Closeout/FinalCloseout/FinalCloseout'));

export class CustomerApp extends Component {
  constructor(props) {
    super(props);

    this.state = { hasError: false, error: undefined, info: undefined };
  }

  componentDidMount() {
    const { loadUser, initOnboarding, loadInternalSchema } = this.props;

    loadInternalSchema();
    loadUser();
    initOnboarding();
  }

  componentDidCatch(error, info) {
    this.setState({
      hasError: true,
      error,
      info,
    });
  }

  render() {
    const props = this.props;
    const { userIsLoggedIn, loginIsLoading } = props;
    const { hasError } = this.state;

    return (
      <>
        <div className="my-move site" id="app-root">
          <CUIHeader />
          <BypassBlock />
          <GovBanner />

          {userIsLoggedIn ? <CustomerLoggedInHeader /> : <LoggedOutHeader />}

          <main role="main" className="site__content my-move-container" id="main">
            <ConnectedLogoutOnInactivity />

            <div className="usa-grid">
              {props.swaggerError && (
                <div className="grid-container">
                  <div className="grid-row">
                    <div className="grid-col-12">
                      <Alert type="error" heading="An error occurred">
                        There was an error contacting the server.
                      </Alert>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {hasError && <SomethingWentWrong />}

            {!userIsLoggedIn && !loginIsLoading && (
              <Routes>
                <Route path={generalRoutes.SIGN_IN_PATH} element={<SignIn />} />
                <Route path={generalRoutes.PRIVACY_SECURITY_POLICY_PATH} element={<PrivacyPolicyStatement />} />
                <Route path={generalRoutes.ACCESSIBILITY_PATH} element={<AccessibilityStatement />} />
                <Route
                  path="*"
                  element={
                    (loginIsLoading && <LoadingPlaceholder />) ||
                    (!userIsLoggedIn && <Navigate to="/sign-in" replace />) || <NotFound />
                  }
                />
              </Routes>
            )}
            {!hasError && !props.swaggerError && userIsLoggedIn && (
              <Routes>
                {/* no auth */}
                <Route path={generalRoutes.SIGN_IN_PATH} element={<SignIn />} />
                <Route path={generalRoutes.PRIVACY_SECURITY_POLICY_PATH} element={<PrivacyPolicyStatement />} />
                <Route path={generalRoutes.ACCESSIBILITY_PATH} element={<AccessibilityStatement />} />

                {/* auth required */}
                <Route end path="/ppm" element={<PpmLanding />} />

                {/* ROOT */}
                <Route path={generalRoutes.HOME_PATH} end element={<Home />} />

                {getWorkflowRoutes(props)}

                <Route end path={customerRoutes.SHIPMENT_MOVING_INFO_PATH} element={<MovingInfo />} />
                <Route end path="/moves/:moveId/edit" element={<Edit />} />
                <Route end path={customerRoutes.EDIT_PROFILE_PATH} element={<EditProfile />} />
                <Route end path={customerRoutes.SERVICE_INFO_EDIT_PATH} element={<EditServiceInfo />} />
                <Route path={customerRoutes.SHIPMENT_CREATE_PATH} element={<ConnectedCreateOrEditMtoShipment />} />
                <Route end path={customerRoutes.PROFILE_PATH} element={<Profile />} />
                <Route end path={customerRoutes.SHIPMENT_EDIT_PATH} element={<ConnectedCreateOrEditMtoShipment />} />
                <Route path={customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH} element={<EstimatedWeightsProGear />} />
                <Route
                  end
                  path={customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH}
                  element={<EstimatedIncentive />}
                />
                <Route end path={customerRoutes.SHIPMENT_PPM_ADVANCES_PATH} element={<Advance />} />
                <Route end path={customerRoutes.CONTACT_INFO_EDIT_PATH} element={<EditContactInfo />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_ABOUT_PATH} element={<About />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH} element={<WeightTickets />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH} element={<WeightTickets />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_REVIEW_PATH} element={<PPMReview />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_EXPENSES_PATH} element={<Expenses />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH} element={<Expenses />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_COMPLETE_PATH} element={<PPMFinalCloseout />} />
                <Route path={customerRoutes.ORDERS_EDIT_PATH} element={<EditOrders />} />
                <Route path={customerRoutes.ORDERS_AMEND_PATH} element={<AmendOrders />} />
                <Route path="/moves/:moveId/review/edit-date-and-location" element={<EditDateAndLocation />} />
                <Route path="/moves/:moveId/review/edit-weight" element={<EditWeight />} />
                <Route end path="/weight-ticket-examples" element={<WeightTicketExamples />} />
                <Route end path="/trailer-criteria" element={<TrailerCriteria />} />
                <Route end path="/allowable-expenses" element={<AllowableExpenses />} />
                <Route end path="/infected-upload" element={<InfectedUpload />} />
                <Route end path="/processing-upload" element={<ProcessingUpload />} />
                <Route path="/moves/:moveId/ppm-payment-request-intro" element={<PPMPaymentRequestIntro />} />
                <Route path="/moves/:moveId/ppm-weight-ticket" element={<WeightTicket />} />
                <Route path="/moves/:moveId/ppm-expenses-intro" element={<ExpensesLanding />} />
                <Route path="/moves/:moveId/ppm-expenses" element={<ExpensesUpload />} />
                <Route path="/moves/:moveId/ppm-payment-review" element={<PaymentReview />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_PRO_GEAR_PATH} element={<ProGear />} />
                <Route end path={customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH} element={<ProGear />} />
                <Route end path="/ppm-customer-agreement" element={<CustomerAgreementLegalese />} />

                {/* Errors */}
                <Route
                  end
                  path="/forbidden"
                  element={
                    <div className="usa-grid">
                      <h2>You are forbidden to use this endpoint</h2>
                    </div>
                  }
                />
                <Route
                  end
                  path="/server_error"
                  element={
                    <div className="usa-grid">
                      <h2>We are experiencing an internal server error</h2>
                    </div>
                  }
                />
                <Route end path="/invalid-permissions" element={<InvalidPermissions />} />

                {/* 404 */}
                <Route
                  path="*"
                  element={
                    (loginIsLoading && <LoadingPlaceholder />) ||
                    (!userIsLoggedIn && <Navigate to="/sign-in" replace />) || <NotFound />
                  }
                />
              </Routes>
            )}
          </main>
          <Footer />
        </div>
        <div id="modal-root"></div>
      </>
    );
  }
}

CustomerApp.propTypes = {
  loadInternalSchema: PropTypes.func,
  loadUser: PropTypes.func,
  initOnboarding: PropTypes.func,
  loginIsLoading: PropTypes.bool,
  userIsLoggedIn: PropTypes.bool,
  conusStatus: PropTypes.string,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      hhgFlow: PropTypes.bool,
      ghcFlow: PropTypes.bool,
    }),
  }).isRequired,
};

CustomerApp.defaultProps = {
  loadInternalSchema: no_op,
  loadUser: no_op,
  initOnboarding: no_op,
  loginIsLoading: false,
  userIsLoggedIn: false,
  conusStatus: '',
  context: {
    flags: {
      hhgFlow: false,
      ghcFlow: false,
    },
  },
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const move = selectCurrentMove(state) || {};

  return {
    loginIsLoading: selectGetCurrentUserIsLoading(state),
    userIsLoggedIn: selectIsLoggedIn(state),
    currentServiceMemberId: serviceMemberId,
    lastMoveIsCanceled: selectHasCanceledMove(state),
    moveId: move?.id,
    conusStatus: selectConusStatus(state),
    swaggerError: state.swaggerInternal.hasErrored,
  };
};
const mapDispatchToProps = {
  loadInternalSchema,
  loadUser: loadUserAction,
  initOnboarding: initOnboardingAction,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(CustomerApp));
