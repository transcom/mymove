import React from 'react';
import { Route } from 'react-router-dom-old';
import { every, some, get, findKey, pick } from 'lodash';

import { generalRoutes, customerRoutes } from 'constants/routes';
import CustomerPrivateRoute from 'containers/CustomerPrivateRoute/CustomerPrivateRoute';
import WizardPage from 'shared/WizardPage';
import generatePath from 'shared/WizardPage/generatePath';
import { no_op } from 'shared/utils';
import { NULL_UUID, SHIPMENT_OPTIONS, CONUS_STATUS } from 'shared/constants';
import BackupContact from 'pages/MyMove/Profile/BackupContact';
import ProfileReview from 'scenes/Review/ProfileReview';

import Home from 'pages/MyMove/Home';
import ConusOrNot from 'pages/MyMove/ConusOrNot';
import DodInfo from 'pages/MyMove/Profile/DodInfo';
import SMName from 'pages/MyMove/Profile/Name';
import DutyLocation from 'pages/MyMove/Profile/DutyLocation';
import ContactInfo from 'pages/MyMove/Profile/ContactInfo';
import Orders from 'pages/MyMove/Orders';
import UploadOrders from 'pages/MyMove/UploadOrders';
import SelectShipmentType from 'pages/MyMove/SelectShipmentType';
import PpmDateAndLocations from 'scenes/Moves/Ppm/DateAndLocation';
import PpmWeight from 'scenes/Moves/Ppm/Weight';
import BackupMailingAddress from 'pages/MyMove/Profile/BackupMailingAddress';
import ResidentialAddress from 'pages/MyMove/Profile/ResidentialAddress';
import Review from 'pages/MyMove/Review/Review';
import Agreement from 'pages/MyMove/Agreement';

const PageNotInFlow = ({ location }) => (
  <div className="usa-grid">
    <h1>Missing Context</h1>
    You are trying to load a page that the system does not have context for. Please go to the home page and try again.
  </div>
);

// USE THESE FOR STUBBING OUT FUTURE WORK
// const Placeholder = props => {
//   return (
//     <WizardPage
//       handleSubmit={() => undefined}
//       pageList={props.pageList}
//       pageKey={props.pageKey}
//     >
//       <div className="Todo-phase2">
//         <h1>Placeholder for {props.title}</h1>
//         <h2>{props.description}</h2>
//       </div>
//     </WizardPage>
//   );
// };

// const stub = (key, pages, description) => ({ match }) => (
//   <Placeholder
//     pageList={pages}
//     pageKey={key}
//     title={key}
//     description={description}
//   />
// );

const always = () => true;
const never = () => false;
// Todo: update this when moves can be completed
const myFirstRodeo = (props) => !props.lastMoveIsCanceled;
const notMyFirstRodeo = (props) => props.lastMoveIsCanceled;
const hasPPM = ({ move }) => {
  return Boolean(move?.mtoShipments?.some((mtoShipment) => mtoShipment.shipmentType === SHIPMENT_OPTIONS.PPM));
};
const inGhcFlow = (props) => props.context.flags.ghcFlow;
const isCurrentMoveSubmitted = ({ move }) => {
  return get(move, 'status', 'DRAFT') === 'SUBMITTED';
};

const pages = {
  [customerRoutes.CONUS_OCONUS_PATH]: {
    isInFlow: inGhcFlow,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.rank, sm.edipi, sm.affiliation]),
    render:
      (key, pages, description, props) =>
      ({ match }) => {
        return (
          <WizardPage
            handleSubmit={no_op}
            pageList={pages}
            pageKey={key}
            match={match}
            canMoveNext={props.conusStatus === CONUS_STATUS.CONUS}
          >
            <ConusOrNot conusStatus={props.conusStatus} />
          </WizardPage>
        );
      },
  },
  [customerRoutes.DOD_INFO_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.rank, sm.edipi, sm.affiliation]),
    render:
      () =>
      ({ history }) =>
        <DodInfo push={history.push} />,
  },
  [customerRoutes.NAME_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.first_name, sm.last_name]),
    render:
      () =>
      ({ history }) =>
        <SMName push={history.push} />,
  },
  [customerRoutes.CONTACT_INFO_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) =>
      sm.is_profile_complete ||
      (every([sm.telephone, sm.personal_email]) && some([sm.phone_is_preferred, sm.email_is_preferred])),
    render:
      () =>
      ({ history }) =>
        <ContactInfo push={history.push} />,
  },
  [customerRoutes.CURRENT_DUTY_LOCATION_PATH]: {
    isInFlow: myFirstRodeo,

    // api for duty location always returns an object, even when duty location is not set
    // if there is no duty location, that object will have a null uuid
    isComplete: ({ sm }) => sm.is_profile_complete || get(sm, 'current_location.id', NULL_UUID) !== NULL_UUID,
    render:
      () =>
      ({ history }) =>
        <DutyLocation push={history.push} />,
    description: 'current duty location',
  },
  [customerRoutes.CURRENT_ADDRESS_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || Boolean(sm.residential_address),
    render:
      () =>
      ({ history }) =>
        <ResidentialAddress push={history.push} />,
  },
  [customerRoutes.BACKUP_ADDRESS_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || Boolean(sm.backup_mailing_address),
    render:
      () =>
      ({ history }) =>
        <BackupMailingAddress push={history.push} />,
  },
  [customerRoutes.BACKUP_CONTACTS_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm, orders, move, ppm, backupContacts }) => {
      return sm.is_profile_complete || backupContacts.length > 0;
    },
    render:
      () =>
      ({ history }) =>
        <BackupContact push={history.push} />,
    description: 'Backup contacts',
  },
  [generalRoutes.HOME_PATH]: {
    isInFlow: (props) => {
      return myFirstRodeo(props) && inGhcFlow(props);
    },
    isComplete: never,
    render:
      (key, pages) =>
      ({ history }) => {
        return <Home history={history} />;
      },
  },
  '/profile-review': {
    isInFlow: notMyFirstRodeo,
    isComplete: always,
    render:
      (key, pages) =>
      ({ match }) =>
        <ProfileReview pages={pages} pageKey={key} match={match} />,
  },
  [customerRoutes.ORDERS_INFO_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders }) =>
      every([
        orders.orders_type,
        orders.issue_date,
        orders.report_by_date,
        get(orders, 'new_duty_location.id', NULL_UUID) !== NULL_UUID,
      ]),
    render:
      (key, pages) =>
      ({ history }) =>
        <Orders push={history.push} />,
  },
  [customerRoutes.ORDERS_UPLOAD_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, uploads }) =>
      get(orders, 'uploaded_orders.uploads', []).length > 0 || uploads.length > 0,
    render:
      (key, pages, description, props) =>
      ({ history }) =>
        <UploadOrders push={history.push} />,
    description: 'Upload your orders',
  },
  [customerRoutes.SHIPMENT_SELECT_TYPE_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, move }) => get(move, 'mtoShipments', []).length > 0,
    render:
      () =>
      ({ history }) =>
        <SelectShipmentType push={history.push} />,
  },
  '/moves/:moveId/ppm-start': {
    isInFlow: hasPPM,
    isComplete: ({ sm, orders, move, ppm }) => {
      return ppm && every([ppm.original_move_date, ppm.pickup_postal_code, ppm.destination_postal_code]);
    },
    render:
      (key, pages) =>
      ({ match }) =>
        <PpmDateAndLocations pages={pages} pageKey={key} match={match} />,
  },
  '/moves/:moveId/ppm-incentive': {
    isInFlow: hasPPM,
    isComplete: ({ sm, orders, move, ppm }) =>
      get(ppm, 'weight_estimate', null) && get(ppm, 'weight_estimate', 0) !== 0,
    render:
      (key, pages) =>
      ({ match }) =>
        <PpmWeight pages={pages} pageKey={key} match={match} />,
  },
  [customerRoutes.MOVE_REVIEW_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, move, ppm, mtoShipment }) => isCurrentMoveSubmitted(move),
    render: () => (props) => {
      return <Review {...props} push={props.history.push} />;
    },
  },
  [customerRoutes.MOVE_AGREEMENT_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, move, ppm, mtoShipment }) => isCurrentMoveSubmitted(move),
    render: () => () => <Agreement />,
  },
};

export const getPagesInFlow = ({ move, conusStatus, lastMoveIsCanceled, context }) =>
  Object.keys(pages).filter((pageKey) => {
    const page = pages[pageKey];
    return page.isInFlow({ move, conusStatus, lastMoveIsCanceled, context });
  });

export const getNextIncompletePage = ({
  conusStatus = '',
  lastMoveIsCanceled = false,
  serviceMember = {},
  orders = {},
  uploads = [],
  move = {},
  ppm = {},
  mtoShipment = {},
  backupContacts = [],
  context = {},
  excludeHomePage = false,
}) => {
  excludeHomePage && delete pages['/'];
  const rawPath = findKey(
    pages,
    (p) =>
      p.isInFlow({ move, conusStatus, lastMoveIsCanceled, context }) &&
      !p.isComplete({ sm: serviceMember, orders, uploads, move, ppm, mtoShipment, backupContacts }),
  );
  const compiledPath = generatePath(rawPath, {
    serviceMemberId: get(serviceMember, 'id'),
    moveId: get(move, 'id'),
  });
  return compiledPath;
};

export const getWorkflowRoutes = (props) => {
  const flowProps = pick(props, ['move', 'conusStatus', 'lastMoveIsCanceled', 'context']);
  const pageList = getPagesInFlow(flowProps);
  return Object.keys(pages).map((key) => {
    const currPage = pages[key];
    if (currPage.isInFlow(flowProps)) {
      const render = currPage.render(key, pageList, currPage.description, props);
      return <CustomerPrivateRoute exact path={key} key={key} render={render} />;
    } else {
      return <Route exact path={key} key={key} component={PageNotInFlow} />;
    }
  });
};
