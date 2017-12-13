const fs = require('fs');
const path = require('path');
const webpack = require('webpack');

function deleteJS(dir) {
    var files = [];
    if (fs.existsSync(dir)) {
        files = fs.readdirSync(dir);
        files.forEach(function (file, index) {
            var curDir = dir + "/" + file;
            if (fs.statSync(curDir).isDirectory()) { // recurse  
                deleteJS(curDir);
            } else { // delete file  
                fs.unlinkSync(curDir);
            }
        });
        fs.rmdirSync(dir);
    }
};

//总是清除之前生成的文件
var jsDir = path.resolve(__dirname, './static/js');
deleteJS(jsDir);

module.exports = {
    entry: {
        vendor: ['vue', 'vue-router', 'iview-style', 'iview', 'lodash', 'axios', 'moment'],
        app: ['./ts/main.ts']
    },
    output: {
        filename: '[name].min.js',
        path: jsDir,
        publicPath: '/static/js/'
    },
    resolve: {
        extensions: ['.ts', '.js', '.vue'],
        alias: {
            'iview-style': 'iview/dist/styles/iview.css'
        }
    },
    watch: true,
    module: {
        rules: [
            {
                test: /\.css$/,
                loader: 'style-loader!css-loader'
            },
            {
                test: /\.vue$/,
                loader: 'vue-loader',
                options: {
                    esModule: true
                }
            },
            {
                test: /\.ts$/,
                loader: 'ts-loader',
                exclude: /node_modules/,
                options: {
                    appendTsSuffixTo: [/\.vue$/]
                }
            },
            {
                test: /\.(png|jpg|gif|woff|woff2|svg|eot|ttf)$/,
                loader: 'url-loader',
                options: {
                    name: './static/img/[hash].[ext]'
                }
            }
        ]
    },
    plugins: [
        new webpack.optimize.CommonsChunkPlugin({
            name: 'vendor',
            minChunks: Infinity
        })
    ]
}

