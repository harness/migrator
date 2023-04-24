// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Harness Upgrade',
  tagline: 'Upgrade from First Gen to Next Gen',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://your-docusaurus-test-site.com',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/migrator/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'harness', // Usually your GitHub org/user name.
  projectName: 'migrator', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: '/',
          sidebarPath: require.resolve('./sidebars.js'),
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/harness.png',
      colorMode: {
        disableSwitch: true,
      },
      navbar: {
        title: 'Harness Upgrade',
        logo: {
          alt: 'Harness logo',
          src: 'img/harness.png',
        },
        items: [
          {
            href: 'https://developer.harness.io/docs/platform/variables-and-expressions/harness-variables/#migrating-firstgen-expressions-to-nextgen',
            label: 'Expressions',
            position: 'right',
          },
          {
            href: 'https://developer.harness.io/docs/continuous-delivery/get-started/upgrading/upgrade-nextgen-cd/',
            label: 'Why upgrade?',
            position: 'right',
          },
          {
            href: 'https://github.com/harness/migrator',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      // footer: {
      //   // style: 'dark',
      //   // copyright: `Copyright © ${new Date().getFullYear()} Harness, Inc.`,
      // },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
  plugins: [
    [
      require.resolve("@cmfcmf/docusaurus-search-local"),
      {
        // Options here
      },
    ],
  ],
  clientModules: [
    require.resolve('./global.js'),
  ]
};

module.exports = config;
