import React from 'react';
import PrivateRoute from 'shared/User/PrivateRoute';
import WizardPage from 'shared/WizardPage';

import Agreement from 'scenes/Legalese';
import Transition from 'scenes/Moves/Transition';
import MoveType from 'scenes/Moves/MoveTypeWizard';
import PpmSize from 'scenes/Moves/Ppm/PPMSizeWizard';
import PpmWeight from 'scenes/Moves/Ppm/Weight';

const Placeholder = props => {
  return (
    <WizardPage
      handleSubmit={() => undefined}
      pageList={props.pageList}
      pageKey={props.pageKey}
    >
      <div className="Todo">
        <h1>Placeholder for {props.title}</h1>
        <h2>{props.description}</h2>
      </div>
    </WizardPage>
  );
};

const stub = (key, pages, description) => ({ match }) => (
  <Placeholder
    pageList={pages}
    pageKey={key}
    title={key}
    description={description}
  />
);
const goHome = props => () => props.push('/');

const incompleteServiceMember = props => !props.hasCompleteProfile;
const hasMove = props => props.hasMove;
const hasHHG = ({ hasMove, selectedMoveType }) =>
  hasMove && selectedMoveType !== 'PPM';
const hasPPM = ({ hasMove, selectedMoveType }) =>
  hasMove && selectedMoveType !== 'HHG';
const isCombo = ({ hasMove, selectedMoveType }) =>
  hasMove && selectedMoveType === 'COMBO';
const pages = {
  '/service-member/:id/create': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Create your profile',
  },
  '/service-member/:id/name': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Name',
  },
  '/service-member/:id/contact-info': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Your contact info',
  },
  '/service-member/:id/duty-station': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'current duty station',
  },
  '/service-member/:id/residence-address': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Current residence address',
  },
  '/service-member/:id/backup-mailing-address': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Backup mailing address',
  },
  '/service-member/:id/backup-contacts': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Backup contacts',
  },
  '/service-member/:id/transition': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: "OK, your profile's complete",
  },
  '/orders/:id/': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Tell us about your move orders',
  },
  '/orders/:id/upload': {
    isInFlow: incompleteServiceMember,
    render: (key, pages, description, props) => ({ match }) => {
      return (
        <WizardPage handleSubmit={goHome(props)} pageList={pages} pageKey={key}>
          <div className="Todo">
            <h1>Placeholder for {key}</h1>
            <h2>Upload your orders</h2>
          </div>
        </WizardPage>
      );
    },
  },
  '/moves/:moveId': {
    isInFlow: hasMove,
    render: (key, pages) => ({ match }) => (
      <MoveType pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/schedule': {
    isInFlow: hasHHG,
    render: stub,
    description: 'Pick a move date',
  },
  '/moves/:moveId/address': {
    isInFlow: hasHHG,
    render: stub,
    description: 'enter your addresses',
  },

  '/moves/:moveId/ppm-transition': {
    isInFlow: isCombo,
    render: (key, pages) => ({ match }) => (
      <WizardPage handleSubmit={() => undefined} pageList={pages} pageKey={key}>
        <Transition />
      </WizardPage>
    ),
  },
  '/moves/:moveId/ppm-start': {
    isInFlow: state => state.selectedMoveType === 'PPM',
    render: stub,
    description: 'pickup zip, destination zip, secondary pickup, temp storage',
  },
  '/moves/:moveId/ppm-size': {
    isInFlow: hasPPM,
    render: (key, pages) => ({ match }) => (
      <PpmSize pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/ppm-incentive': {
    isInFlow: hasPPM,
    render: (key, pages) => ({ match }) => (
      <PpmWeight pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/review': {
    isInFlow: hasMove,
    render: stub,
    description: 'Review',
  },
  '/moves/:moveId/agreement': {
    isInFlow: hasMove,
    render: (key, pages, description, props) => ({ match }) => {
      return (
        <WizardPage handleSubmit={goHome(props)} pageList={pages} pageKey={key}>
          <Agreement match={match} />
        </WizardPage>
      );
    },
  },
};
export const getPageList = state =>
  Object.keys(pages).filter(pageKey => {
    const page = pages[pageKey];
    return page.isInFlow(state);
  });

export const getWorkflowRoutes = props => {
  const pageList = getPageList(props);
  return pageList.map(key => {
    const currPage = pages[key];
    console.log(key, currPage.description, props);
    const render = currPage.render(key, pageList, currPage.description, props);
    return <PrivateRoute exact path={key} key={key} render={render} />;
  });
};
