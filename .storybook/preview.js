// eslint-disable-next-line import/no-extraneous-dependencies
import 'happo-plugin-storybook/register';

import '../src/index.scss';
import '../src/ghc_index.scss';

import '../src/icons';

export const parameters = {
  options: {
    storySort: {
      order: ['Global', 'Components', 'Office Components', 'Customer Components'],
    },
  },
  a11y: {
    // axe-core configurationOptions (https://github.com/dequelabs/axe-core/blob/develop/doc/API.md#parameters-1)
    config: {},
    // axe-core optionsParameter (https://github.com/dequelabs/axe-core/blob/develop/doc/API.md#options-parameter)
    options: {},
  },
};
