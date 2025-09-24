module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'subject-max-length': [
        2,
        'always',
        120
    ],
    'header-max-length': [2, "always",70]
  }
}

