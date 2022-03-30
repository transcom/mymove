import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { func, node, number, string } from 'prop-types';
import { generatePath } from 'react-router';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './MovingInfo.module.scss';

import { customerRoutes, generalRoutes } from 'constants/routes';
import ScrollToTop from 'components/ScrollToTop';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { fetchLatestOrders as fetchLatestOrdersAction } from 'shared/Entities/modules/orders';
import { formatWeight } from 'utils/formatters';
import { selectCurrentOrders, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { RouteProps } from 'types/router';

const IconSection = ({ icon, headline, children }) => (
  <Grid row className={styles.IconSection}>
    <Grid col="auto" className={styles.SectionIcon}>
      <FontAwesomeIcon size="lg" icon={icon} />
    </Grid>
    <Grid col="fill" className={styles.SectionContent}>
      <h2 className={styles.SectionHeadline}>{headline} </h2>
      {children}
    </Grid>
  </Grid>
);

IconSection.propTypes = {
  icon: string.isRequired,
  headline: string.isRequired,
  children: node.isRequired,
};

export class MovingInfo extends Component {
  componentDidMount() {
    const { serviceMemberId, fetchLatestOrders } = this.props;
    fetchLatestOrders(serviceMemberId);
  }

  render() {
    const { entitlementWeight, history, match } = this.props;

    const nextPath = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
      moveId: match.params.moveId,
    });

    return (
      <GridContainer>
        <ScrollToTop />
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <h1 className={styles.ShipmentsHeader}>Things to know about selecting shipments</h1>
            <SectionWrapper className={styles.Wrapper}>
              <IconSection
                icon="weight-hanging"
                headline={`You can move ${formatWeight(entitlementWeight)} in this move.`}
              >
                <p>You&apos;ll have to pay for any excess weight the government moves.</p>
              </IconSection>
              <IconSection icon="pencil-alt" headline="You don't need to get the details perfect.">
                <p>
                  After you submit this information, you&apos;ll talk to a move counselor. They will verify your choices
                  and help identify more complicated situations.
                </p>
                <p>
                  If you use movers, they will be your point of contact throughout your move and can also help you make
                  changes to your shipments.
                </p>
              </IconSection>
              <IconSection icon="truck-moving" headline="If you use movers, they will:">
                <div className={styles.IconSectionList}>
                  <ul>
                    <li>Help estimate how much your belongings weigh</li>
                    <li>Set pack and pickup dates based on your preferred pickup date</li>
                    <li>Contact you after you talk to a move counselor</li>
                    <li>Be your main point of contact during your move</li>
                  </ul>
                </div>
              </IconSection>
              <IconSection icon="car" headline="It's common to move some things yourself.">
                <p>
                  Most people have professionals move most of their things, but handle a few important things
                  themselves.
                </p>
              </IconSection>
              <IconSection icon="hand-holding-usd" headline="You can get paid for things you move yourself.">
                <p>
                  The government will pay you for moving belongings that you document by weight. (This is a PPM, or
                  DITY.)
                </p>
              </IconSection>
            </SectionWrapper>

            <WizardNavigation
              isFirstPage
              showFinishLater
              onNextClick={() => {
                history.push(nextPath);
              }}
              onCancelClick={() => {
                history.push(generalRoutes.HOME_PATH);
              }}
            />
          </Grid>
        </Grid>
      </GridContainer>
    );
  }
}

MovingInfo.propTypes = {
  entitlementWeight: number,
  fetchLatestOrders: func.isRequired,
  serviceMemberId: string.isRequired,
  ...RouteProps,
};

MovingInfo.defaultProps = {
  entitlementWeight: 0,
};

function mapStateToProps(state) {
  const orders = selectCurrentOrders(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const entitlementWeight = orders.has_dependents
    ? serviceMember?.weight_allotment?.total_weight_self_plus_dependents
    : serviceMember?.weight_allotment?.total_weight_self;
  const serviceMemberId = serviceMember?.id;

  return {
    entitlementWeight,
    serviceMemberId,
  };
}

const mapDispatchToProps = {
  fetchLatestOrders: fetchLatestOrdersAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(MovingInfo);
